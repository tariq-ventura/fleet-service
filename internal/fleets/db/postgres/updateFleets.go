package fleets_db_postgres

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) UpdateFleets(id uuid.UUID, updates map[string]any) (*fleets_domain.Fleet, *interfaces.Error) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":       "fleets",
		"db.operation":  "update",
		"db.type":       "postgresql",
		"equipments.Id": id,
	})
	defer operationSpan.End()

	result, err := pc.ListFleetsById(id)

	if err != nil {
		return nil, err
	}

	update := pc.client.WithContext(spanCtx).Model(&result).Where("id = ?", id).Updates(updates)

	if update.Error != nil {
		switch {
		case errors.Is(update.Error, gorm.ErrDuplicatedKey):
			return nil, &interfaces.Error{
				Error:      "fleet_already_exists",
				Message:    "Ya existe una flota con ese codigo",
				StatusCode: http.StatusConflict,
			}

		default:
			pc.logging.LogError("fleet_database_error", map[string]any{"error": update.Error.Error()})
			return nil, &interfaces.Error{
				Error:      "database_error",
				Message:    "No se pudo actualizar la flota",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	if update.RowsAffected == 0 {
		return nil, &interfaces.Error{
			Error:      "fleet_not_found",
			Message:    "La flota no existe",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return result, nil
}
