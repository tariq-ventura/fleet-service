package equipments_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (eh *EquipmentHanlder) ListEquipmentsById(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "El identificador no es un UUID válido",
		})
		return
	}

	span, _ := eh.trace.StartSpan(
		ctx,
		"equipments.create_equipment",
		map[string]any{
			"http.method":    "GET",
			"http.route":     "/api/v1/equipments/${id}",
			"http.params.id": id,
		},
	)
	defer span.End()

	dbSpan, dbCtx := eh.trace.StartSpan(ctx, "equipments.database.connection", map[string]any{
		"db.name": "equipments",
	})
	database := eh.db
	dbSpan.End()

	operationSpan, _ := eh.trace.StartSpan(dbCtx, "equipments.database.operations", map[string]any{
		"db.name":      "equipments",
		"db.operation": "list",
		"db.params.id": id,
	})
	defer operationSpan.End()

	result, erro := database.ListEquipmentsById(id)

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
