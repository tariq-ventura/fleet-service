package equipments_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	equipments_dto "github.com/tariq-ventura/fleet-service/internal/equipments/dto"
)

func (eh *EquipmentHanlder) CreateEquipment(c *gin.Context) {
	var input equipments_dto.CreateEquipmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Los datos enviados no son válidos", "detail": err.Error()})
		return
	}

	vehicle := &equipments_domain.Equipment{
		Description: strings.TrimSpace(input.Description), Status: strings.TrimSpace(input.Status),
		Type: strings.TrimSpace(input.Type), Year: input.Year, Color: strings.TrimSpace(input.Color),
		Brand: strings.TrimSpace(input.Brand), Model: strings.TrimSpace(input.Model), Group: strings.TrimSpace(input.Group),
		Tags: strings.TrimSpace(input.Tags), Driver: strings.TrimSpace(input.Driver), RemoteID: strings.TrimSpace(input.RemoteID),
	}
	if responseError := eh.db.CreateEquipment(vehicle, c.Request.Context()); responseError != nil {
		c.JSON(responseError.StatusCode, gin.H{"error": responseError.Error, "message": responseError.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vehículo registrado correctamente", "data": vehicle})
}
