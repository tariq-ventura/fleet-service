package equipments_db_postgres

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	equipments_dto "github.com/tariq-ventura/fleet-service/internal/equipments/dto"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvalidStatusTransition = errors.New(
		"invalid equipment status transition",
	)

	ErrConcurrentStatusChange = errors.New(
		"equipment status changed concurrently",
	)
)

func (pc *PostgresClient) UpdateEquipmentStatus(equipmentId uuid.UUID, newStatus equipments_domain.EquipmentStatus) (*equipments_domain.EquipmentStatusHistory, *equipments_domain.Equipment, *interfaces.Error) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":       "equipments",
		"db.operation":  "update",
		"db.type":       "postgresql",
		"equipments.Id": equipmentId,
	})
	defer operationSpan.End()

	var request equipments_dto.UpdateEquipmentStatusRequest
	var createdHistory equipments_domain.EquipmentStatusHistory
	var updatedEquipment equipments_domain.Equipment

	transactionError := pc.client.
		WithContext(spanCtx).
		Transaction(func(tx *gorm.DB) error {
			var equipment equipments_domain.Equipment

			result := tx.
				Clauses(clause.Locking{
					Strength: "UPDATE",
				}).
				First(&equipment, "id = ?", equipmentId)

			if result.Error != nil {
				return result.Error
			}

			if equipment.Status == newStatus {
				return fmt.Errorf(
					"%w: equipment is already %s",
					ErrInvalidStatusTransition,
					newStatus,
				)
			}

			if !equipment.Status.CanTransitionTo(newStatus) {
				return fmt.Errorf(
					"%w: %s -> %s",
					ErrInvalidStatusTransition,
					equipment.Status,
					newStatus,
				)
			}

			previousStatus := equipment.Status

			// Se agrega también el estado actual al WHERE.
			// Es una protección adicional ante cambios concurrentes.
			updateResult := tx.
				Model(&equipments_domain.Equipment{}).
				Where(
					"id = ? AND status = ?",
					equipmentId,
					previousStatus,
				).
				Updates(map[string]any{
					"status": newStatus,
				})

			if updateResult.Error != nil {
				return updateResult.Error
			}

			if updateResult.RowsAffected == 0 {
				return ErrConcurrentStatusChange
			}

			createdHistory = equipments_domain.EquipmentStatusHistory{
				ID:          uuid.New(),
				EquipmentID: equipmentId,
				FromStatus:  previousStatus,
				ToStatus:    newStatus,
				Reason:      strings.TrimSpace(request.Reason),
				ChangedAt:   time.Now().UTC(),
			}

			if err := tx.Create(&createdHistory).Error; err != nil {
				return err
			}

			if err := tx.
				First(&updatedEquipment, "id = ?", equipmentId).
				Error; err != nil {

				return err
			}

			return nil
		})

	if transactionError != nil {
		switch {
		case errors.Is(
			transactionError,
			gorm.ErrRecordNotFound,
		):
			return nil, nil, &interfaces.Error{
				Error:      "equipment_not_found",
				Message:    "La maquinaria no existe",
				StatusCode: http.StatusNotFound,
			}

		case errors.Is(
			transactionError,
			ErrInvalidStatusTransition,
		):
			return nil, nil, &interfaces.Error{
				Error:      "invalid_status_transition",
				Message:    transactionError.Error(),
				StatusCode: http.StatusUnprocessableEntity,
			}

		case errors.Is(
			transactionError,
			ErrConcurrentStatusChange,
		):
			return nil, nil, &interfaces.Error{
				Error: "concurrent_status_change",
				Message: "El estado de la maquinaria cambió " +
					"durante la operación. Consulte nuevamente.",
				StatusCode: http.StatusConflict,
			}

		default:
			return nil, nil, &interfaces.Error{
				Error:      "database_error",
				Message:    "No se pudo cambiar el estado",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	return &createdHistory, &updatedEquipment, nil
}
