package equipments_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (eh *EquipmentHanlder) RemoveEquipment(c *gin.Context) {
	id, ok := validations.ParseUUIDParameter(c, "id")
	if !ok {
		return
	}
	if responseError := eh.db.RemoveEquipment(id, c.Request.Context()); responseError != nil {
		c.JSON(responseError.StatusCode, gin.H{"error": responseError.Error, "message": responseError.Message})
		return
	}
	c.Status(http.StatusNoContent)
}
