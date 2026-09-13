package geofences_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	geofences_domain "github.com/tariq-ventura/fleet-service/internal/geofences/domain"
	geofences_dto "github.com/tariq-ventura/fleet-service/internal/geofences/dto"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (h *GeofenceHandler) CreateGeofence(c *gin.Context) {
	var input geofences_dto.CreateGeofenceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Los datos enviados no son válidos", "detail": err.Error()})
		return
	}
	item := &geofences_domain.Geofence{GeofenceID: strings.TrimSpace(input.GeofenceID), Name: strings.TrimSpace(input.Name), Group: strings.TrimSpace(input.Group), AdditionalMargin: input.AdditionalMargin, Latitude: input.Latitude, Longitude: input.Longitude}
	if err := h.db.CreateGeofence(item, c.Request.Context()); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Geocerca registrada correctamente", "data": item})
}

func (h *GeofenceHandler) ListGeofences(c *gin.Context) {
	page := validations.ParsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := validations.ParsePositiveInt(c.DefaultQuery("pageSize", "20"), 20)
	items, err, total := h.db.ListGeofences(page, pageSize, strings.TrimSpace(c.Query("search")), c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "pagination": gin.H{"page": page, "pageSize": pageSize, "total": total, "totalPages": validations.CalculateTotalPages(total, pageSize)}})
}

func (h *GeofenceHandler) ListGeofenceById(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	item, err := h.db.ListGeofenceById(id, c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *GeofenceHandler) UpdateGeofence(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	var input geofences_dto.UpdateGeofenceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Los datos enviados no son válidos", "detail": err.Error()})
		return
	}
	updates := make(map[string]any)
	if input.GeofenceID != nil {
		updates["geofence_id"] = strings.TrimSpace(*input.GeofenceID)
	}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.Group != nil {
		updates["group"] = strings.TrimSpace(*input.Group)
	}
	if input.AdditionalMargin != nil {
		updates["additional_margin"] = *input.AdditionalMargin
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
	item, err := h.db.UpdateGeofence(id, updates, c.Request.Context())
	if err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Geocerca actualizada correctamente", "data": item})
}

func (h *GeofenceHandler) RemoveGeofence(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	if err := h.db.RemoveGeofence(id, c.Request.Context()); err != nil {
		c.JSON(err.StatusCode, gin.H{"error": err.Error, "message": err.Message})
		return
	}
	c.Status(http.StatusNoContent)
}
