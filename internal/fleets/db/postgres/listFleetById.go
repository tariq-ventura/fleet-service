package fleets_db_postgres

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) ListFleetsById(id uuid.UUID) (*fleets_domain.Fleet, *interfaces.Error) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":       "fleets",
		"db.operation":  "list",
		"db.type":       "postgresql",
		"equipments.Id": id,
	})
	defer operationSpan.End()

	var fleet fleets_domain.Fleet

	result := pc.client.WithContext(spanCtx).First(&fleet, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &interfaces.Error{
			Error:      "equipment_not_found",
			Message:    "La flota no existe",
			StatusCode: http.StatusNotFound,
		}
	}

	if result.Error != nil {
		pc.logging.LogError("fleet_database_error", map[string]any{"error": result.Error.Error()})
		return nil, &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo consultar la flota",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &fleet, nil
}
