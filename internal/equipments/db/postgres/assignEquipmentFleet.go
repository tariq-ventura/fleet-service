package equipments_db_postgres

import (
	"errors"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/validations"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (pc *PostgresClient) AssignEquipmentFleet(fleetId, equipmentId uuid.UUID) (*equipments_domain.Equipment, *interfaces.Error) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":       "equipments",
		"db.operation":  "update",
		"db.type":       "postgresql",
		"equipments.Id": equipmentId,
		"fleets.Id":     fleetId,
	})
	defer operationSpan.End()

	var updatedEquipment equipments_domain.Equipment

	err := pc.client.WithContext(spanCtx).Transaction(func(tx *gorm.DB) error {
		var fleet fleets_domain.Fleet

		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).First(&fleet, "id = ?", fleetId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return validations.ErrFleetNotFound
			}

			return err
		}

		var equipment equipments_domain.Equipment

		if err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
		}).First(&equipment, "id = ?", equipmentId).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return validations.ErrEquipmentNotFound
			}

			return err
		}

		if equipment.FleetID != nil &&
			*equipment.FleetID == fleetId {

			updatedEquipment = equipment
			return nil
		}

		if equipment.FleetID != nil {
			return validations.ErrEquipmentAlreadyAssigned
		}

		if !validations.CanChangeFleet(equipment.Status) {
			return validations.ErrEquipmentInOperation
		}
		result := tx.
			Model(&equipments_domain.Equipment{}).
			Where(
				"id = ? AND fleet_id IS NULL",
				equipmentId,
			).
			Update("fleet_id", fleetId)

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return validations.ErrEquipmentAlreadyAssigned
		}

		return tx.
			First(
				&updatedEquipment,
				"id = ?",
				equipmentId,
			).
			Error
	})

	if err != nil {
		pc.logging.LogError("database_error", map[string]any{"error": err.Error()})
		return nil, validations.HandleFleetMembershipError(err)
	}

	return &updatedEquipment, nil
}
