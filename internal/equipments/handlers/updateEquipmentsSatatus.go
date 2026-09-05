package equipments_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	equipments_dto "github.com/tariq-ventura/fleet-service/internal/equipments/dto"
)

func (eh *EquipmentHanlder) UpdateEquipmentStatus(c *gin.Context) {
	ctx := c.Request.Context()
	equipmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "El identificador no es un UUID válido",
		})
		return
	}

	var request equipments_dto.UpdateEquipmentStatusRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		eh.logs.LogWarning("invalid_request", map[string]any{"details": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Los datos enviados no son válidos",
		})
		return
	}

	newStatus := equipments_domain.EquipmentStatus(
		strings.ToUpper(strings.TrimSpace(request.Status)),
	)

	if !newStatus.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_status",
			"message": "El estado indicado no es válido",
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
		"db.operation": "update",
	})
	defer operationSpan.End()

	updatedEquipment, createdHistory, dberr := database.UpdateEquipmentStatus(equipmentID, newStatus)

	if dberr != nil {
		c.JSON(dberr.StatusCode, gin.H{
			"error":   dberr.Error,
			"message": dberr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Estado actualizado correctamente",
		"data": gin.H{
			"equipment":  updatedEquipment,
			"transition": createdHistory,
		},
	})
}
