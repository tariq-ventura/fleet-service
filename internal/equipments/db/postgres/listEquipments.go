package equipments_db_postgres

import (
	"context"
	"net/http"
	"strings"

	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
)

func (pc *PostgresClient) ListEquipments(page, pageSize int, equipmentType, status, brand, search string, ctx context.Context) ([]equipments_domain.Equipment, *interfaces.Error, int64) {
	if pageSize > 100 {
		pageSize = 100
	}

	query := pc.client.WithContext(ctx).Model(&equipments_domain.Equipment{})
	if equipmentType != "" {
		query = query.Where("LOWER(type) = ?", strings.ToLower(strings.TrimSpace(equipmentType)))
	}
	if status != "" {
		query = query.Where("LOWER(status) = ?", strings.ToLower(strings.TrimSpace(status)))
	}
	if brand != "" {
		query = query.Where("LOWER(brand) = ?", strings.ToLower(strings.TrimSpace(brand)))
	}
	if search != "" {
		pattern := "%" + strings.TrimSpace(search) + "%"
		query = query.Where("description ILIKE ? OR type ILIKE ? OR brand ILIKE ? OR model ILIKE ? OR remote_id ILIKE ? OR tags ILIKE ?", pattern, pattern, pattern, pattern, pattern, pattern)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, databaseError("No se pudo contar los vehículos", err, pc), 0
	}

	vehicles := make([]equipments_domain.Equipment, 0)
	result := query.Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&vehicles)
	if result.Error != nil {
		return nil, databaseError("No se pudieron consultar los vehículos", result.Error, pc), 0
	}

	return vehicles, nil, total
}

func databaseError(message string, err error, pc *PostgresClient) *interfaces.Error {
	pc.logging.LogError("database_error", map[string]any{"error": err.Error()})
	return &interfaces.Error{Error: "database_error", Message: message, StatusCode: http.StatusInternalServerError}
}
