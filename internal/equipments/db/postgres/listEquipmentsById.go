package equipments_db_postgres

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) ListEquipmentsById(id uuid.UUID) (*equipments_domain.Equipment, *interfaces.Error) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":       "equipments",
		"db.operation":  "list",
		"db.type":       "postgresql",
		"equipments.Id": id,
	})
	defer operationSpan.End()

	var equipment equipments_domain.Equipment

	result := pc.client.WithContext(spanCtx).First(&equipment, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, &interfaces.Error{
			Error:      "equipment_not_found",
			Message:    "La maquinaria no existe",
			StatusCode: http.StatusNotFound,
		}
	}

	if result.Error != nil {
		pc.logging.LogError("database_error", map[string]any{"error": result.Error.Error()})
		return nil, &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo consultar la maquinaria",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &equipment, nil
}
