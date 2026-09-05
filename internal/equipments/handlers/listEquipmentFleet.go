package equipments_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (eh *EquipmentHanlder) ListEquipmentsFleets(c *gin.Context) {
	ctx := c.Request.Context()
	fleetID, ok := validations.ParseUUIDParameter(c, "fleetID")
	if !ok {
		return
	}

	span, _ := eh.trace.StartSpan(
		ctx,
		"equipments.create_equipment",
		map[string]any{
			"http.method":         "GET",
			"http.route":          "/api/v1/equipments/fleets/${fleetID}",
			"http.params.fleetID": fleetID,
		},
	)
	defer span.End()

	dbSpan, dbCtx := eh.trace.StartSpan(ctx, "equipments.database.connection", map[string]any{
		"db.name": "equipments",
	})
	database := eh.db
	dbSpan.End()

	operationSpan, _ := eh.trace.StartSpan(dbCtx, "fleets.database.operations", map[string]any{
		"db.name":           "fleets",
		"db.operation":      "list",
		"db.params.fleetID": fleetID,
	})
	defer operationSpan.End()

	equipments, fleet, erro := database.ListEquipmentsFleets(fleetID)

	if erro != nil {
		c.JSON(erro.StatusCode, gin.H{
			"error":   erro.Error,
			"message": erro.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"fleet": fleet,
		"data":  equipments,
	})
}
