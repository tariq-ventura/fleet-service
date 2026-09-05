package fleets_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (fh *FleetHanlder) ListFleetsById(c *gin.Context) {
	ctx := c.Request.Context()
	id, ok := validations.ParseUUIDParameter(c, "fleetID")
	if !ok {
		return
	}

	span, _ := fh.trace.StartSpan(
		ctx,
		"equipments.create_equipment",
		map[string]any{
			"http.method":    "GET",
			"http.route":     "/api/v1/fleets/${id}",
			"http.params.id": id,
		},
	)
	defer span.End()

	dbSpan, dbCtx := fh.trace.StartSpan(ctx, "fleets.database.connection", map[string]any{
		"db.name": "fleets",
	})
	database := fh.db
	dbSpan.End()

	operationSpan, _ := fh.trace.StartSpan(dbCtx, "fleets.database.operations", map[string]any{
		"db.name":      "fleets",
		"db.operation": "list",
		"db.params.id": id,
	})
	defer operationSpan.End()

	result, erro := database.ListFleetsById(id)

	if erro != nil {
		c.JSON(erro.StatusCode, gin.H{
			"error":   erro.Error,
			"message": erro.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
