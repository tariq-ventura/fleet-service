package equipments_db_postgres

import (
	"errors"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/validations"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (pc *PostgresClient) RemoveEquipmentFleet(fleetId, equipmentId uuid.UUID) *interfaces.Error {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":       "equipments",
		"db.operation":  "update",
		"db.type":       "postgresql",
		"equipments.Id": equipmentId,
		"fleets.Id":     fleetId,
	})
	defer operationSpan.End()

	err := pc.client.WithContext(spanCtx).Transaction(func(tx *gorm.DB) error {
		var equipment equipments_domain.Equipment

		if err := tx.
			Clauses(clause.Locking{
				Strength: "UPDATE",
			}).
			First(&equipment, "id = ?", equipmentId).
			Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return validations.ErrEquipmentNotFound
			}

			return err
		}

		if equipment.FleetID == nil ||
			*equipment.FleetID != fleetId {

			return validations.ErrEquipmentNotInFleet
		}

		if !validations.CanChangeFleet(equipment.Status) {
			return validations.ErrEquipmentInOperation
		}

		result := tx.
			Model(&equipments_domain.Equipment{}).
			Where(
				"id = ? AND fleet_id = ?",
				equipmentId,
				fleetId,
			).
			Update("fleet_id", nil)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return validations.ErrEquipmentNotInFleet
		}

		return nil
	})

	if err != nil {
		return validations.HandleFleetMembershipError(err)
	}

	return nil
}
