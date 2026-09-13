package router

import (
	"github.com/gin-gonic/gin"
	equipments_handlers "github.com/tariq-ventura/fleet-service/internal/equipments/handlers"
)

func (ro *Routes) EquipmentsRoutes(r *gin.Engine) {
	er := equipments_handlers.NewEquipmentHanlder(ro.Context, ro.EquipmentsDB, ro.Trace, ro.Logging)
	routes := r.Group("/api/v1/equipments")
	{
		routes.POST("", er.CreateEquipment)
		routes.GET("", er.ListEquipments)
		routes.GET("/:id", er.ListEquipmentsById)
		routes.PATCH("/:id", er.UpdateEquipments)
		routes.DELETE("/:id", er.RemoveEquipment)
		routes.PATCH("/:id/status", er.UpdateEquipmentStatus)
		routes.GET("/:id/status-history", er.ListEquipmentStatusHistory)
	}
}

// VehiclesRoutes is an alias that exposes the Startrack terminology while
// keeping /equipments available for the existing n8n flow.
func (ro *Routes) VehiclesRoutes(r *gin.Engine) {
	er := equipments_handlers.NewEquipmentHanlder(ro.Context, ro.EquipmentsDB, ro.Trace, ro.Logging)
	routes := r.Group("/api/v1/vehicles")
	{
		routes.POST("", er.CreateEquipment)
		routes.GET("", er.ListEquipments)
		routes.GET("/:id", er.ListEquipmentsById)
		routes.PATCH("/:id", er.UpdateEquipments)
		routes.DELETE("/:id", er.RemoveEquipment)
		routes.PATCH("/:id/status", er.UpdateEquipmentStatus)
		routes.GET("/:id/status-history", er.ListEquipmentStatusHistory)
	}
}
