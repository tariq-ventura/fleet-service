package router

import (
	"github.com/gin-gonic/gin"
	equipments_handlers "github.com/tariq-ventura/fleet-service/internal/equipments/handlers"
	fleets_handlers "github.com/tariq-ventura/fleet-service/internal/fleets/handlers"
)

func (ro *Routes) FleetsRoutes(r *gin.Engine) {
	fr := fleets_handlers.NewFleetHanlder(ro.Context, ro.FleetsDB, ro.Trace, ro.Logging)
	er := equipments_handlers.NewEquipmentHanlder(ro.Context, ro.EquipmentsDB, ro.Trace, ro.Logging)

	routes := r.Group("/api/v1/fleets")
	{
		routes.POST("", fr.CreateFleets)
		routes.GET("", fr.ListFleets)
		routes.GET("/:fleetID", fr.ListFleetsById)
		routes.PATCH("/:fleetID", fr.UpdateFleets)

		routes.PUT("/:fleetID/equipments/:equipmentID", er.AssignEquipmentFleet)
		routes.DELETE("/:fleetID/equipments/:equipmentID", er.RemoveEquipmentFleet)

		routes.GET("/:fleetID/equipments", er.ListEquipmentsFleets)
	}
}
