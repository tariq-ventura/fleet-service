package equipments_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (eh *EquipmentHanlder) ListEquipmentStatusHistory(c *gin.Context) {
	ctx := c.Request.Context()
	equipmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "El identificador no es un UUID válido",
		})
		return
	}

	dbSpan, dbCtx := eh.trace.StartSpan(ctx, "equipments.database.connection", map[string]any{
		"db.name": "equipments",
	})
	database := eh.db
	dbSpan.End()

	operationSpan, _ := eh.trace.StartSpan(dbCtx, "equipments.database.operations", map[string]any{
		"db.name":      "equipments",
		"db.operation": "list",
	})
	defer operationSpan.End()

	history, dberr := database.ListEquipmentStatusHistory(equipmentID)

	if dberr != nil {
		c.JSON(dberr.StatusCode, gin.H{
			"error":   dberr.Error,
			"message": dberr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": history,
	})
}
