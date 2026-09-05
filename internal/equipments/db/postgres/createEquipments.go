package equipments_db_postgres

import (
	"errors"
	"net/http"

	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) CreateEquipment(data equipments_domain.Equipment) *interfaces.Error {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":              "equipments",
		"db.operation":         "insert",
		"db.type":              "postgresql",
		"equipments.code":      data.Code,
		"equipments.fleet":     data.FleetID,
		"equipments.createdAt": data.CreatedAt,
	})
	defer operationSpan.End()

	result := pc.client.WithContext(spanCtx).Create(&data)

	if result.Error != nil {
		switch {
		case errors.Is(result.Error, gorm.ErrDuplicatedKey):
			pc.logging.LogWarning("equipment_already_exists", nil)
			return &interfaces.Error{
				Error:      "equipment_already_exists",
				Message:    "Ya existe una maquinaria con ese código o número de serie",
				StatusCode: http.StatusConflict,
			}

		case errors.Is(result.Error, gorm.ErrForeignKeyViolated):
			pc.logging.LogWarning("fleet_not_found", nil)
			return &interfaces.Error{
				Error:      "fleet_not_found",
				Message:    "La flota especificada no existe",
				StatusCode: http.StatusBadRequest,
			}

		default:
			pc.logging.LogError("database_error", map[string]any{"error": result.Error})
			return &interfaces.Error{
				Error:      "database_error",
				Message:    "No se pudo registrar la maquinaria",
				StatusCode: http.StatusInternalServerError,
			}
		}
	}

	pc.logging.LogInfo("PostgreSQL insert success", map[string]interface{}{"insertedID": data.ID})
	return nil
}
