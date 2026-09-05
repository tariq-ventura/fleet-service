package equipments_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (eh *EquipmentHanlder) ListEquipments(c *gin.Context) {
	ctx := c.Request.Context()
	equipmentType := c.Query("type")
	status := c.Query("status")
	brand := c.Query("brand")
	search := strings.TrimSpace(c.Query("search"))

	span, _ := eh.trace.StartSpan(
		ctx,
		"equipments.list_equipments",
		map[string]any{
			"http.method":               "GET",
			"http.route":                "/api/v1/equipments",
			"http.params.equipmentType": equipmentType,
			"http.params.status":        status,
			"http.params.brand":         brand,
			"http.params.search":        search,
		},
	)
	defer span.End()

	page := validations.ParsePositiveInt(c.DefaultQuery("page", "1"), 1)
	pageSize := validations.ParsePositiveInt(c.DefaultQuery("pageSize", "20"), 20)

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

	result, err, total := database.ListEquipments(page, pageSize, equipmentType, status, brand, search)

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
