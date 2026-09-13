package equipments_db_postgres

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) UpdateEquipments(id uuid.UUID, updates map[string]any, ctx context.Context) (*equipments_domain.Equipment, *interfaces.Error) {
	if _, responseError := pc.ListEquipmentsById(id, ctx); responseError != nil {
		return nil, responseError
	}

	result := pc.client.WithContext(ctx).Model(&equipments_domain.Equipment{}).Where("id = ?", id).Updates(updates)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, &interfaces.Error{Error: "vehicle_already_exists", Message: "Ya existe un vehículo con esa descripción o ID remoto", StatusCode: http.StatusConflict}
	}
	if result.Error != nil {
		return nil, databaseError("No se pudo actualizar el vehículo", result.Error, pc)
	}

	return pc.ListEquipmentsById(id, ctx)
}
