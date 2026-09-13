package router

import (
	"github.com/gin-gonic/gin"
	geofences_handlers "github.com/tariq-ventura/fleet-service/internal/geofences/handlers"
)

func (ro *Routes) GeofencesRoutes(r *gin.Engine) {
	h := geofences_handlers.NewGeofenceHandler(ro.Context, ro.GeofencesDB, ro.Trace, ro.Logging)
	routes := r.Group("/api/v1/geofences")
	{
		routes.POST("", h.CreateGeofence)
		routes.GET("", h.ListGeofences)
		routes.GET("/:id", h.ListGeofenceById)
		routes.PATCH("/:id", h.UpdateGeofence)
		routes.DELETE("/:id", h.RemoveGeofence)
	}
}
