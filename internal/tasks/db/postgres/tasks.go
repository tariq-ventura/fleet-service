package tasks_db_postgres

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	tasks_domain "github.com/tariq-ventura/fleet-service/internal/tasks/domain"
	"gorm.io/gorm"
)

func (pc *PostgresClient) responseError(message string, err error) *interfaces.Error {
	pc.logging.LogError("task_database_error", map[string]any{"error": err.Error()})
	return &interfaces.Error{Error: "database_error", Message: message, StatusCode: http.StatusInternalServerError}
}
func (pc *PostgresClient) CreateTask(data *tasks_domain.Task, ctx context.Context) *interfaces.Error {
	result := pc.client.WithContext(ctx).Create(data)
	if result.Error == nil {
		return nil
	}
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return &interfaces.Error{Error: "task_already_exists", Message: "Ya existe una tarea con ese identificador", StatusCode: http.StatusConflict}
	}
	return pc.responseError("No se pudo registrar la tarea", result.Error)
}
func (pc *PostgresClient) ListTasks(page, pageSize int, status, taskType, search string, ctx context.Context) ([]tasks_domain.Task, *interfaces.Error, int64) {
	if pageSize > 100 {
		pageSize = 100
	}
	query := pc.client.WithContext(ctx).Model(&tasks_domain.Task{})
	if status != "" {
		query = query.Where("LOWER(status) = ?", strings.ToLower(strings.TrimSpace(status)))
	}
	if taskType != "" {
		query = query.Where("LOWER(type) = ?", strings.ToLower(strings.TrimSpace(taskType)))
	}
	if search != "" {
		pattern := "%" + strings.TrimSpace(search) + "%"
		query = query.Where("task_id ILIKE ? OR title ILIKE ? OR description ILIKE ? OR origin ILIKE ? OR destination ILIKE ? OR assignee ILIKE ?", pattern, pattern, pattern, pattern, pattern, pattern)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, pc.responseError("No se pudo contar las tareas", err), 0
	}
	items := make([]tasks_domain.Task, 0)
	if err := query.Order("scheduled_date DESC").Limit(pageSize).Offset((page - 1) * pageSize).Find(&items).Error; err != nil {
		return nil, pc.responseError("No se pudieron consultar las tareas", err), 0
	}
	return items, nil, total
}
func (pc *PostgresClient) ListTaskById(id uuid.UUID, ctx context.Context) (*tasks_domain.Task, *interfaces.Error) {
	var item tasks_domain.Task
	err := pc.client.WithContext(ctx).First(&item, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, &interfaces.Error{Error: "task_not_found", Message: "La tarea no existe", StatusCode: http.StatusNotFound}
	}
	if err != nil {
		return nil, pc.responseError("No se pudo consultar la tarea", err)
	}
	return &item, nil
}
func (pc *PostgresClient) UpdateTask(id uuid.UUID, updates map[string]any, ctx context.Context) (*tasks_domain.Task, *interfaces.Error) {
	if _, err := pc.ListTaskById(id, ctx); err != nil {
		return nil, err
	}
	result := pc.client.WithContext(ctx).Model(&tasks_domain.Task{}).Where("id = ?", id).Updates(updates)
	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, &interfaces.Error{Error: "task_already_exists", Message: "Ya existe una tarea con ese identificador", StatusCode: http.StatusConflict}
	}
	if result.Error != nil {
		return nil, pc.responseError("No se pudo actualizar la tarea", result.Error)
	}
	return pc.ListTaskById(id, ctx)
}
func (pc *PostgresClient) RemoveTask(id uuid.UUID, ctx context.Context) *interfaces.Error {
	result := pc.client.WithContext(ctx).Delete(&tasks_domain.Task{}, "id = ?", id)
	if result.Error != nil {
		return pc.responseError("No se pudo eliminar la tarea", result.Error)
	}
	if result.RowsAffected == 0 {
		return &interfaces.Error{Error: "task_not_found", Message: "La tarea no existe", StatusCode: http.StatusNotFound}
	}
	return nil
}
