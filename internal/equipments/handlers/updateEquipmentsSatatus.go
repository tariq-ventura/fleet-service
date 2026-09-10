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

	span, spanCtx := eh.trace.StartSpan(
		ctx,
		"equipments.create_equipment",
		map[string]any{
			"http.method":    "PATCH",
			"http.route":     "/api/v1/equipments/${id}/status",
			"http.params.id": equipmentID,
		},
	)
	defer span.End()

	var request equipments_dto.UpdateEquipmentStatusRequest

	bindSpan, bindCtx := eh.trace.StartSpan(spanCtx, "equipments.create_equipment.BindJson", nil)
	bindError := c.ShouldBindJSON(&request)
	bindSpan.End()

	if bindError != nil {
		eh.logs.LogWarning("invalid request", map[string]interface{}{"error": bindError.Error()})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Los datos enviados no son válidos",
			"detail":  err.Error(),
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

	dbSpan, dbCtx := eh.trace.StartSpan(bindCtx, "equipments.database.connection", map[string]any{
		"db.name": "equipments",
	})
	database := eh.db
	dbSpan.End()

	operationSpan, opCtx := eh.trace.StartSpan(dbCtx, "equipments.database.operations", map[string]any{
		"db.name":      "equipments",
		"db.operation": "update",
	})
	defer operationSpan.End()

	updatedEquipment, createdHistory, dberr := database.UpdateEquipmentStatus(equipmentID, newStatus, opCtx)

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
