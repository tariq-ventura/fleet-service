package router

import (
	"github.com/gin-gonic/gin"
	maintenance_handlers "github.com/tariq-ventura/fleet-service/internal/maintenance/handlers"
)

func (ro *Routes) MaintenanceRoutes(r *gin.Engine) {
	h := maintenance_handlers.NewMaintenanceHandler(ro.Context, ro.MaintenanceDB, ro.Trace, ro.Logging)
	routes := r.Group("/api/v1/maintenance")
	{
		routes.POST("", h.CreateMaintenance)
		routes.GET("", h.ListMaintenance)
		routes.GET("/:id", h.ListMaintenanceById)
		routes.PATCH("/:id", h.UpdateMaintenance)
		routes.DELETE("/:id", h.RemoveMaintenance)
	}
}
