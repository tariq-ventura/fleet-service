package fleets_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	fleets_dto "github.com/tariq-ventura/fleet-service/internal/fleets/dto"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (fh *FleetHanlder) UpdateFleets(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := validations.ParseUUIDParameter(c, "fleetID")
	if !ok {
		return
	}

	span, _ := fh.trace.StartSpan(
		ctx,
		"equipments.create_equipment",
		map[string]any{
			"http.method":    "PATCH",
			"http.route":     "/api/v1/fleets/${id}",
			"http.params.id": id,
		},
	)
	defer span.End()

	var request fleets_dto.UpdateFleetRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		fh.logs.LogWarning("invalid_request", map[string]any{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Los datos enviados no son válidos",
		})
		return
	}

	updates := make(map[string]any)

	if request.Code != nil {
		updates["code"] = strings.ToUpper(
			strings.TrimSpace(*request.Code),
		)
	}

	if request.Name != nil {
		updates["name"] = strings.TrimSpace(*request.Name)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "empty_update",
			"message": "Debe enviar al menos un campo",
		})
		return
	}

	dbSpan, dbCtx := fh.trace.StartSpan(ctx, "fleets.database.connection", map[string]any{
		"db.name": "fleets",
	})
	database := fh.db
	dbSpan.End()

	operationSpan, _ := fh.trace.StartSpan(dbCtx, "fleets.database.operations", map[string]any{
		"db.name":      "fleets",
		"db.operation": "update",
	})
	defer operationSpan.End()

	fleet, erro := database.UpdateFleets(id, updates)

	if erro != nil {
		c.JSON(erro.StatusCode, gin.H{
			"error":   erro.Error,
			"message": erro.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Flota actualizada correctamente",
		"data":    fleet,
	})
}
