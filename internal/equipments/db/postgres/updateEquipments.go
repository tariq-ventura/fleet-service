package equipments_db_postgres

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) UpdateEquipments(id uuid.UUID, updates map[string]any) (*equipments_domain.Equipment, *interfaces.Error) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":       "equipments",
		"db.operation":  "update",
		"db.type":       "postgresql",
		"equipments.Id": id,
	})
	defer operationSpan.End()

	result, err := pc.ListEquipmentsById(id)

	if err != nil {
		return nil, err
	}

	update := pc.client.WithContext(spanCtx).Model(&result).Where("id = ?", id).Updates(updates)

	if update.Error != nil {
		switch {
		case errors.Is(update.Error, gorm.ErrDuplicatedKey):
			return nil, &interfaces.Error{
				Error:      "equipment_already_exists",
				Message:    "Ya existe una maquinaria con ese número de serie",
				StatusCode: http.StatusConflict,
			}

		case errors.Is(update.Error, gorm.ErrForeignKeyViolated):
			return nil, &interfaces.Error{
				Error:      "fleet_not_found",
				Message:    "La flota especificada no existe",
				StatusCode: http.StatusBadRequest,
			}

		default:
			pc.logging.LogError("database_error", map[string]any{"error": update.Error.Error()})
			return nil, &interfaces.Error{
				Error:      "database_error",
				Message:    "No se pudo actualizar la maquinaria",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}
	return result, nil
}
