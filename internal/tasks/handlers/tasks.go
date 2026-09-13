package tasks_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	tasks_domain "github.com/tariq-ventura/fleet-service/internal/tasks/domain"
	tasks_dto "github.com/tariq-ventura/fleet-service/internal/tasks/dto"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var input tasks_dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Los datos enviados no son válidos", "detail": err.Error()})
		return
	}
	item := &tasks_domain.Task{TaskID: strings.TrimSpace(input.TaskID), Title: strings.TrimSpace(input.Title), Description: strings.TrimSpace(input.Description), Type: strings.TrimSpace(input.Type), Status: strings.TrimSpace(input.Status), ScheduledDate: input.ScheduledDate.UTC(), Origin: strings.TrimSpace(input.Origin), Destination: strings.TrimSpace(input.Destination), Latitude: input.Latitude, Longitude: input.Longitude, Assignee: strings.TrimSpace(input.Assignee)}
	if err := h.db.CreateTask(item, c.Request.Context()); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Tarea registrada correctamente", "data": item})
}
func (h *TaskHandler) ListTasks(c *gin.Context) {
	page := validations.ParsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := validations.ParsePositiveInt(c.DefaultQuery("pageSize", "20"), 20)
	items, err, total := h.db.ListTasks(page, pageSize, c.Query("status"), c.Query("type"), strings.TrimSpace(c.Query("search")), c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "pagination": gin.H{"page": page, "pageSize": pageSize, "total": total, "totalPages": validations.CalculateTotalPages(total, pageSize)}})
}
func (h *TaskHandler) ListTaskById(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	item, err := h.db.ListTaskById(id, c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	var input tasks_dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Los datos enviados no son válidos", "detail": err.Error()})
		return
	}
	updates := make(map[string]any)
	add := func(k string, v *string) {
		if v != nil {
			updates[k] = strings.TrimSpace(*v)
		}
	}
	add("task_id", input.TaskID)
	add("title", input.Title)
	add("description", input.Description)
	add("type", input.Type)
	add("status", input.Status)
	add("origin", input.Origin)
	add("destination", input.Destination)
	add("assignee", input.Assignee)
	if input.ScheduledDate != nil {
		updates["scheduled_date"] = input.ScheduledDate.UTC()
	}
	if input.Latitude != nil {
		updates["latitude"] = *input.Latitude
	}
	if input.Longitude != nil {
		updates["longitude"] = *input.Longitude
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty_update", "message": "Debe enviar al menos un campo"})
		return
	}
	item, err := h.db.UpdateTask(id, updates, c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tarea actualizada correctamente", "data": item})
}
func (h *TaskHandler) RemoveTask(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	if err := h.db.RemoveTask(id, c.Request.Context()); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.Status(http.StatusNoContent)
}
