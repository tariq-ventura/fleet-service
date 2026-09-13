package equipments_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	equipments_dto "github.com/tariq-ventura/fleet-service/internal/equipments/dto"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (eh *EquipmentHanlder) UpdateEquipments(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}

	var input equipments_dto.UpdateEquipmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": "Los datos enviados no son válidos", "detail": err.Error()})
		return
	}

	updates := make(map[string]any)
	addTrimmedString := func(column string, value *string) {
		if value != nil {
			updates[column] = strings.TrimSpace(*value)
		}
	}
	addTrimmedString("description", input.Description)
	addTrimmedString("type", input.Type)
	addTrimmedString("color", input.Color)
	addTrimmedString("brand", input.Brand)
	addTrimmedString("model", input.Model)
	addTrimmedString("group", input.Group)
	addTrimmedString("tags", input.Tags)
	addTrimmedString("driver", input.Driver)
	addTrimmedString("remote_id", input.RemoteID)
	if input.Year != nil {
		updates["year"] = *input.Year
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty_update", "message": "Debe enviar al menos un campo"})
		return
	}

	vehicle, responseError := eh.db.UpdateEquipments(id, updates, c.Request.Context())
	if responseError != nil {
		c.JSON(responseError.StatusCode, gin.H{"error": responseError.Error, "message": responseError.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vehículo actualizado correctamente", "data": vehicle})
}
