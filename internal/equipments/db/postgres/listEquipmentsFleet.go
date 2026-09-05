package equipments_db_postgres

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) ListEquipmentsFleets(fleeId uuid.UUID) ([]equipments_domain.Equipment, *fleets_domain.Fleet, *interfaces.Error) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":            "equipments",
		"db.operation":       "list",
		"db.type":            "postgresql",
		"equipments.FleetID": fleeId,
	})
	defer operationSpan.End()

	var fleet fleets_domain.Fleet

	if err := pc.client.WithContext(spanCtx).First(&fleet, "id = ?", fleeId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, &interfaces.Error{
				Error:      "fleet_not_found",
				Message:    "La flota no existe",
				StatusCode: http.StatusNotFound,
			}
		}

		pc.logging.LogError("database_error", map[string]any{"error": err.Error()})
		return nil, nil, &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo consultar la flota.",
			StatusCode: http.StatusInternalServerError,
		}
	}

	var equipment []equipments_domain.Equipment

	result := pc.client.WithContext(spanCtx).Where("fleet_id = ?", fleeId).Order("code ASC").Find(&equipment)

	if result.Error != nil {
		pc.logging.LogError("database_error", map[string]any{"error": result.Error.Error()})
		return nil, nil, &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo consultar la maquinaria",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return equipment, &fleet, nil
}
