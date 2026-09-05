package equipments_db_postgres

import (
	"net/http"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
)

func (pc *PostgresClient) ListEquipmentStatusHistory(id uuid.UUID) ([]equipments_domain.EquipmentStatusHistory, *interfaces.Error) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":       "equipments",
		"db.operation":  "list",
		"db.type":       "postgresql",
		"equipments.Id": id,
	})
	defer operationSpan.End()

	_, err := pc.ListEquipmentsById(id)

	if err != nil {
		return nil, err
	}

	var history []equipments_domain.EquipmentStatusHistory

	result := pc.client.WithContext(spanCtx).Where("equipment_id = ?", id).Order("changed_at DESC").Find(&history)

	if result.Error != nil {
		pc.logging.LogError("database_error", map[string]any{"error": result.Error.Error()})
		return nil, &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo consultar el historial",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return history, nil
}
