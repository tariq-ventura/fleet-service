package equipments_db_postgres

import (
	"context"
	"errors"
	"net/http"

	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) CreateEquipment(data *equipments_domain.Equipment, ctx context.Context) *interfaces.Error {
	result := pc.client.WithContext(ctx).Create(data)
	if result.Error == nil {
		return nil
	}

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return &interfaces.Error{Error: "vehicle_already_exists", Message: "Ya existe un vehículo con esa descripción o ID remoto", StatusCode: http.StatusConflict}
	}

	pc.logging.LogError("database_error", map[string]any{"error": result.Error.Error()})
	return &interfaces.Error{Error: "database_error", Message: "No se pudo registrar el vehículo", StatusCode: http.StatusInternalServerError}
}
