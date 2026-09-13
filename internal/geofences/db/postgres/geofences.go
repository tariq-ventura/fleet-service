package geofences_db_postgres

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	geofences_domain "github.com/tariq-ventura/fleet-service/internal/geofences/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
)

func (pc *PostgresClient) responseError(message string, err error) *interfaces.Error {
	pc.logging.LogError("geofence_database_error", map[string]any{"error": err.Error()})
	return &interfaces.Error{Error: "database_error", Message: message, StatusCode: http.StatusInternalServerError}
}

func (pc *PostgresClient) CreateGeofence(data *geofences_domain.Geofence, ctx context.Context) *interfaces.Error {
	result := pc.client.WithContext(ctx).Create(data)
	if result.Error == nil {
		return nil
	}
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return &interfaces.Error{Error: "geofence_already_exists", Message: "Ya existe una geocerca con ese identificador", StatusCode: http.StatusConflict}
	}
	return pc.responseError("No se pudo registrar la geocerca", result.Error)
}

func (pc *PostgresClient) ListGeofences(page, pageSize int, search string, ctx context.Context) ([]geofences_domain.Geofence, *interfaces.Error, int64) {
	if pageSize > 100 {
		pageSize = 100
	}
	query := pc.client.WithContext(ctx).Model(&geofences_domain.Geofence{})
	if search != "" {
		pattern := "%" + strings.TrimSpace(search) + "%"
		query = query.Where("geofence_id ILIKE ? OR name ILIKE ? OR \"group\" ILIKE ?", pattern, pattern, pattern)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, pc.responseError("No se pudo contar las geocercas", err), 0
	}
	items := make([]geofences_domain.Geofence, 0)
	if err := query.Order("created_at DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&items).Error; err != nil {
		return nil, pc.responseError("No se pudieron consultar las geocercas", err), 0
	}
	return items, nil, total
}

func (pc *PostgresClient) ListGeofenceById(id uuid.UUID, ctx context.Context) (*geofences_domain.Geofence, *interfaces.Error) {
	var item geofences_domain.Geofence
	err := pc.client.WithContext(ctx).First(&item, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &interfaces.Error{Error: "geofence_not_found", Message: "La geocerca no existe", StatusCode: http.StatusNotFound}
	}
	if err != nil {
		return nil, pc.responseError("No se pudo consultar la geocerca", err)
	}
	return &item, nil
}

func (pc *PostgresClient) UpdateGeofence(id uuid.UUID, updates map[string]any, ctx context.Context) (*geofences_domain.Geofence, *interfaces.Error) {
	if _, err := pc.ListGeofenceById(id, ctx); err != nil {
		return nil, err
	}
	result := pc.client.WithContext(ctx).Model(&geofences_domain.Geofence{}).Where("id = ?", id).Updates(updates)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, &interfaces.Error{Error: "geofence_already_exists", Message: "Ya existe una geocerca con ese identificador", StatusCode: http.StatusConflict}
	}
	if result.Error != nil {
		return nil, pc.responseError("No se pudo actualizar la geocerca", result.Error)
	}
	return pc.ListGeofenceById(id, ctx)
}

func (pc *PostgresClient) RemoveGeofence(id uuid.UUID, ctx context.Context) *interfaces.Error {
	result := pc.client.WithContext(ctx).Delete(&geofences_domain.Geofence{}, "id = ?", id)
	if result.Error != nil {
		return pc.responseError("No se pudo eliminar la geocerca", result.Error)
	}
	if result.RowsAffected == 0 {
		return &interfaces.Error{Error: "geofence_not_found", Message: "La geocerca no existe", StatusCode: http.StatusNotFound}
	}
	return nil
}
