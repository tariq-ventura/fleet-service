package equipments_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	equipments_dto "github.com/tariq-ventura/fleet-service/internal/equipments/dto"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (eh *EquipmentHanlder) UpdateEquipmentStatus(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	var input equipments_dto.UpdateEquipmentStatusRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Los datos enviados no son válidos", "detail": err.Error()})
		return
	}

	history, vehicle, responseError := eh.db.UpdateEquipmentStatus(id, strings.TrimSpace(input.Status), strings.TrimSpace(input.Reason), c.Request.Context())
	if responseError != nil {
		c.JSON(responseError.StatusCode, gin.H{"error": responseError.Error, "message": responseError.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Estado actualizado correctamente", "data": gin.H{"vehicle": vehicle, "transition": history}})
}
