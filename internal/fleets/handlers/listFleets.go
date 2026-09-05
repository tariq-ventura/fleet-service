package fleets_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (fh *FleetHanlder) ListFleets(c *gin.Context) {
	ctx := c.Request.Context()
	search := strings.TrimSpace(c.Query("search"))

	span, _ := fh.trace.StartSpan(
		ctx,
		"equipments.create_equipment",
		map[string]any{
			"http.method":        "GET",
			"http.route":         "/api/v1/fleets",
			"http.params.search": search,
		},
	)
	defer span.End()

	page := validations.ParsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := validations.ParsePositiveInt(c.DefaultQuery("pageSize", "20"), 20)

	dbSpan, dbCtx := fh.trace.StartSpan(ctx, "fleets.database.connection", map[string]any{
		"db.name": "fleets",
	})
	database := fh.db
	dbSpan.End()

	operationSpan, _ := fh.trace.StartSpan(dbCtx, "fleets.database.operations", map[string]any{
		"db.name":      "fleets",
		"db.operation": "list",
	})
	defer operationSpan.End()

	result, err, total := database.ListFleets(page, pageSize, search)

	if err != nil {
		c.JSON(err.StatusCode, gin.H{
			"error":   err.Error,
			"message": err.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
		"pagination": gin.H{
			"page":       page,
			"pageSize":   pageSize,
			"total":      total,
			"totalPages": validations.CalculateTotalPages(total, pageSize),
		},
	})
}
