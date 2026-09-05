package equipments_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	equipments_dto "github.com/tariq-ventura/fleet-service/internal/equipments/dto"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

func (eh *EquipmentHanlder) CreateEquipment(c *gin.Context) {
	ctx := c.Request.Context()
	var request equipments_dto.CreateEquipmentRequest

	span, _ := eh.trace.StartSpan(
		ctx,
		"equipments.create_equipment",
		map[string]any{
			"http.method": "POST",
			"http.route":  "/api/v1/equipments",
		},
	)
	defer span.End()

	bindSpan, _ := eh.trace.StartSpan(ctx, "equipments.create_equipment.BindJson", nil)
	err := c.ShouldBindJSON(&request)
	bindSpan.End()

	if err != nil {
		eh.logs.LogWarning("invalid request", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Los datos enviados no son válidos",
			"detail":  err.Error(),
		})
		return
	}

	fleetSpan, _ := eh.trace.StartSpan(ctx, "equipments.create_equipment.validate_fleetID", nil)
	fleetId := validations.ValidateFleetId(request.FleetID)
	fleetSpan.End()

	if fleetSpan == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_fleet_id",
			"message": "fleetId debe ser un UUID válido",
		})
		return
	}

	equipment := equipments_domain.Equipment{
		ID:           uuid.New(),
		Code:         strings.ToUpper(strings.TrimSpace(request.Code)),
		FleetID:      fleetId,
		Type:         equipments_domain.EquipmentType(request.Type),
		Brand:        strings.TrimSpace(request.Brand),
		Model:        strings.TrimSpace(request.Model),
		SerialNumber: strings.ToUpper(strings.TrimSpace(request.SerialNumber)),
		Year:         request.Year,

		CapacityTons: request.CapacityTons,
		Status:       equipments_domain.StatusAvailable,

		Location: equipments_domain.Location{
			Name:      strings.TrimSpace(request.Location.Name),
			Latitude:  request.Location.Latitude,
			Longitude: request.Location.Longitude,
		},

		EngineHours: request.EngineHours,

		NextMaintenanceHours: request.EngineHours +
			request.MaintenanceInterval,

		FuelPercent: request.FuelPercent,
	}

	dbSpan, dbCtx := eh.trace.StartSpan(ctx, "equipments.database.connection", map[string]any{
		"db.name": "equipments",
	})
	database := eh.db
	dbSpan.End()

	operationSpan, _ := eh.trace.StartSpan(dbCtx, "equipments.database.operations", map[string]any{
		"db.name":              "equipments",
		"db.operation":         "insert",
		"equipments.code":      equipment.Code,
		"equipments.fleet":     equipment.FleetID,
		"equipments.createdAt": equipment.CreatedAt,
	})
	defer operationSpan.End()

	result := database.CreateEquipment(equipment)

	if result != nil {
		c.JSON(result.StatusCode, gin.H{
			"error":   result.Error,
			"message": result.Message,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Maquinaria registrada correctamente",
		"data":    equipment,
	})
}
