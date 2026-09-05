package validations

import (
	"errors"
	"net/http"

	"github.com/tariq-ventura/fleet-service/internal/interfaces"
)

var (
	ErrFleetNotFound = errors.New(
		"fleet not found",
	)

	ErrEquipmentNotFound = errors.New(
		"equipment not found",
	)

	ErrFleetHasEquipment = errors.New(
		"fleet has assigned equipment",
	)

	ErrEquipmentAlreadyAssigned = errors.New(
		"equipment already assigned",
	)

	ErrEquipmentInOperation = errors.New(
		"equipment is in operation",
	)

	ErrEquipmentNotInFleet = errors.New(
		"equipment does not belong to fleet",
	)
)

func HandleFleetMembershipError(err error) *interfaces.Error {
	switch {
	case errors.Is(err, ErrFleetNotFound):
		return &interfaces.Error{
			Error:      "fleet_not_found",
			Message:    "La flota no existe",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, ErrEquipmentNotFound):
		return &interfaces.Error{
			Error:      "equipment_not_found",
			Message:    "La maquinaria no existe",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, ErrEquipmentAlreadyAssigned):
		return &interfaces.Error{
			Error: "equipment_already_assigned",
			Message: "La maquinaria ya pertenece " +
				"a otra flota",
			StatusCode: http.StatusConflict,
		}
	case errors.Is(err, ErrEquipmentNotInFleet):
		return &interfaces.Error{
			Error: "equipment_not_in_fleet",
			Message: "La maquinaria no pertenece " +
				"a esta flota",
			StatusCode: http.StatusNotFound,
		}
	case errors.Is(err, ErrEquipmentInOperation):
		return &interfaces.Error{
			Error: "equipment_in_operation",
			Message: "No se puede mover la maquinaria " +
				"mientras está en operación",
			StatusCode: http.StatusConflict,
		}
	default:
		return &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo completar la operación",
			StatusCode: http.StatusInternalServerError,
		}
	}
}
