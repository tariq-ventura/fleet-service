package maintenance_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	maintenance_domain "github.com/tariq-ventura/fleet-service/internal/maintenance/domain"
	maintenance_dto "github.com/tariq-ventura/fleet-service/internal/maintenance/dto"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (h *MaintenanceHandler) CreateMaintenance(c *gin.Context) {
	var input maintenance_dto.CreateMaintenanceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Los datos enviados no son válidos", "detail": err.Error()})
		return
	}
	item := &maintenance_domain.Maintenance{Vehicle: strings.TrimSpace(input.Vehicle), Reference: strings.TrimSpace(input.Reference), ServiceDate: input.ServiceDate.UTC(), Odometer: input.Odometer, ServiceTime: strings.TrimSpace(input.ServiceTime), HourMeter: input.HourMeter, RepairReason: strings.TrimSpace(input.RepairReason), Provider: strings.TrimSpace(input.Provider), Mechanic: strings.TrimSpace(input.Mechanic), ServiceType: strings.TrimSpace(input.ServiceType)}
	if err := h.db.CreateMaintenance(item, c.Request.Context()); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Mantenimiento registrado correctamente", "data": item})
}
func (h *MaintenanceHandler) ListMaintenance(c *gin.Context) {
	page := validations.ParsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := validations.ParsePositiveInt(c.DefaultQuery("pageSize", "20"), 20)
	items, err, total := h.db.ListMaintenance(page, pageSize, c.Query("vehicle"), c.Query("search"), c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "pagination": gin.H{"page": page, "pageSize": pageSize, "total": total, "totalPages": validations.CalculateTotalPages(total, pageSize)}})
}
func (h *MaintenanceHandler) ListMaintenanceById(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	item, err := h.db.ListMaintenanceById(id, c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}
func (h *MaintenanceHandler) UpdateMaintenance(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	var input maintenance_dto.UpdateMaintenanceRequest
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
	add("vehicle", input.Vehicle)
	add("reference", input.Reference)
	add("service_time", input.ServiceTime)
	add("repair_reason", input.RepairReason)
	add("provider", input.Provider)
	add("mechanic", input.Mechanic)
	add("service_type", input.ServiceType)
	if input.ServiceDate != nil {
		updates["service_date"] = input.ServiceDate.UTC()
	}
	if input.Odometer != nil {
		updates["odometer"] = *input.Odometer
	}
	if input.HourMeter != nil {
		updates["hour_meter"] = *input.HourMeter
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty_update", "message": "Debe enviar al menos un campo"})
		return
	}
	item, err := h.db.UpdateMaintenance(id, updates, c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mantenimiento actualizado correctamente", "data": item})
}
func (h *MaintenanceHandler) RemoveMaintenance(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	if err := h.db.RemoveMaintenance(id, c.Request.Context()); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.Status(http.StatusNoContent)
}
