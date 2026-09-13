package maintenance_db_postgres

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	maintenance_domain "github.com/tariq-ventura/fleet-service/internal/maintenance/domain"
	"gorm.io/gorm"
)

func (pc *PostgresClient) responseError(message string, err error) *interfaces.Error {
	pc.logging.LogError("maintenance_database_error", map[string]any{"error": err.Error()})
	return &interfaces.Error{Error: "database_error", Message: message, StatusCode: http.StatusInternalServerError}
}
func (pc *PostgresClient) CreateMaintenance(data *maintenance_domain.Maintenance, ctx context.Context) *interfaces.Error {
	if err := pc.client.WithContext(ctx).Create(data).Error; err != nil {
		return pc.responseError("No se pudo registrar el mantenimiento", err)
	}
	return nil
}
func (pc *PostgresClient) ListMaintenance(page, pageSize int, vehicle, search string, ctx context.Context) ([]maintenance_domain.Maintenance, *interfaces.Error, int64) {
	if pageSize > 100 {
		pageSize = 100
	}
	query := pc.client.WithContext(ctx).Model(&maintenance_domain.Maintenance{})
	if vehicle != "" {
		query = query.Where("LOWER(vehicle) = ?", strings.ToLower(strings.TrimSpace(vehicle)))
	}
	if search != "" {
		pattern := "%" + strings.TrimSpace(search) + "%"
		query = query.Where("vehicle ILIKE ? OR reference ILIKE ? OR repair_reason ILIKE ? OR provider ILIKE ? OR mechanic ILIKE ? OR service_type ILIKE ?", pattern, pattern, pattern, pattern, pattern, pattern)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, pc.responseError("No se pudo contar los mantenimientos", err), 0
	}
	items := make([]maintenance_domain.Maintenance, 0)
	if err := query.Order("service_date DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&items).Error; err != nil {
		return nil, pc.responseError("No se pudieron consultar los mantenimientos", err), 0
	}
	return items, nil, total
}
func (pc *PostgresClient) ListMaintenanceById(id uuid.UUID, ctx context.Context) (*maintenance_domain.Maintenance, *interfaces.Error) {
	var item maintenance_domain.Maintenance
	err := pc.client.WithContext(ctx).First(&item, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &interfaces.Error{Error: "maintenance_not_found", Message: "El mantenimiento no existe", StatusCode: http.StatusNotFound}
	}
	if err != nil {
		return nil, pc.responseError("No se pudo consultar el mantenimiento", err)
	}
	return &item, nil
}
func (pc *PostgresClient) UpdateMaintenance(id uuid.UUID, updates map[string]any, ctx context.Context) (*maintenance_domain.Maintenance, *interfaces.Error) {
	if _, err := pc.ListMaintenanceById(id, ctx); err != nil {
		return nil, err
	}
	if err := pc.client.WithContext(ctx).Model(&maintenance_domain.Maintenance{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return nil, pc.responseError("No se pudo actualizar el mantenimiento", err)
	}
	return pc.ListMaintenanceById(id, ctx)
}
func (pc *PostgresClient) RemoveMaintenance(id uuid.UUID, ctx context.Context) *interfaces.Error {
	result := pc.client.WithContext(ctx).Delete(&maintenance_domain.Maintenance{}, "id = ?", id)
	if result.Error != nil {
		return pc.responseError("No se pudo eliminar el mantenimiento", result.Error)
	}
	if result.RowsAffected == 0 {
		return &interfaces.Error{Error: "maintenance_not_found", Message: "El mantenimiento no existe", StatusCode: http.StatusNotFound}
	}
	return nil
}
