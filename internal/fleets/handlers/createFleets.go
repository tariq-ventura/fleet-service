package fleets_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	fleets_dto "github.com/tariq-ventura/fleet-service/internal/fleets/dto"
)

func (fh *FleetHanlder) CreateFleets(c *gin.Context) {
	ctx := c.Request.Context()
	var request fleets_dto.CreateFleetRequest

	span, _ := fh.trace.StartSpan(
		ctx,
		"fleets.create_fleet",
		map[string]any{
			"http.method": "POST",
			"http.route":  "/api/v1/fleets",
		},
	)
	defer span.End()

	bindSpan, _ := fh.trace.StartSpan(ctx, "fleets.create_fleets.BindJson", nil)
	err := c.ShouldBindJSON(&request)
	bindSpan.End()

	if err != nil {
		fh.logs.LogWarning("invalid_request", map[string]any{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Los datos enviados no son válidos",
		})
		return
	}

	fleet := fleets_domain.Fleet{
		ID:   uuid.New(),
		Code: strings.ToUpper(strings.TrimSpace(request.Code)),
		Name: strings.TrimSpace(request.Name),
	}

	dbSpan, dbCtx := fh.trace.StartSpan(ctx, "fleets.database.connection", map[string]any{
		"db.name": "fleets",
	})
	database := fh.db
	dbSpan.End()

	operationSpan, _ := fh.trace.StartSpan(dbCtx, "fleets.database.operations", map[string]any{
		"db.name":              "fleets",
		"db.operation":         "insert",
		"fleets.code":      fleet.Code,
		"fleets.createdAt": fleet.CreatedAt,
	})
	defer operationSpan.End()

	result := database.CreateFleets(fleet)

	if result != nil {
		c.JSON(result.StatusCode, gin.H{
			"error":   result.Error,
			"message": result.Message,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Flota creada correctamente",
		"data":    fleet,
	})
}
