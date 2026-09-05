package equipments_handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	equipments_dto "github.com/tariq-ventura/fleet-service/internal/equipments/dto"
)

func (eh *EquipmentHanlder) UpdateEquipments(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_id",
			"message": "El identificador no es un UUID válido",
		})
		return
	}

	var request equipments_dto.UpdateEquipmentRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid_request",
			"message": "Los datos enviados no son válidos",
			"detail":  err.Error(),
		})
		return
	}

	updates := make(map[string]any)

	if request.FleetID != nil {
		fleetID, err := uuid.Parse(*request.FleetID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid_fleet_id",
				"message": "fleetId debe ser un UUID válido",
			})
			return
		}

		updates["fleet_id"] = fleetID
	}

	if request.Type != nil {
		updates["type"] = strings.ToUpper(
			strings.TrimSpace(*request.Type),
		)
	}

	if request.Brand != nil {
		updates["brand"] = strings.TrimSpace(*request.Brand)
	}

	if request.Model != nil {
		updates["model"] = strings.TrimSpace(*request.Model)
	}

	if request.SerialNumber != nil {
		updates["serial_number"] = strings.ToUpper(
			strings.TrimSpace(*request.SerialNumber),
		)
	}

	if request.Year != nil {
		updates["year"] = *request.Year
	}

	if request.CapacityTons != nil {
		updates["capacity_tons"] = *request.CapacityTons
	}

	if request.Location != nil {
		if request.Location.Name != nil {
			updates["location_name"] = strings.TrimSpace(
				*request.Location.Name,
			)
		}

		if request.Location.Latitude != nil {
			updates["latitude"] = *request.Location.Latitude
		}

		if request.Location.Longitude != nil {
			updates["longitude"] = *request.Location.Longitude
		}
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "empty_update",
			"message": "Debe enviar al menos un campo para actualizar",
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

	result, erro := database.UpdateEquipments(id, updates)

	if erro != nil {
		c.JSON(erro.StatusCode, gin.H{
			"error":   erro.Error,
			"message": erro.Message,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Maquinaria actualizada correctamente",
		"data":    result,
	})
}
