package fleets_db_postgres

import (
	"errors"
	"net/http"

	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) CreateFleets(fleet fleets_domain.Fleet) *interfaces.Error {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":          "fleets",
		"db.operation":     "insert",
		"db.type":          "postgresql",
		"fleets.code":      fleet.Code,
		"fleets.createdAt": fleet.CreatedAt,
	})
	defer operationSpan.End()

	result := pc.client.WithContext(spanCtx).Create(&fleet)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return &interfaces.Error{
				Error:      "fleet_already_exists",
				Message:    "Ya existe una flota con ese código",
				StatusCode: http.StatusConflict,
			}
		}
		pc.logging.LogError("fleet_database_error", map[string]any{"error": result.Error.Error()})
		return &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo crear la flota",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return nil
}
