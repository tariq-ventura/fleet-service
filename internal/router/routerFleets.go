package router

import (
	"github.com/gin-gonic/gin"
	fleets_handlers "github.com/tariq-ventura/fleet-service/internal/fleets/handlers"
)

// FleetsRoutes remains as a compatibility endpoint. Startrack's Group value is
// stored directly on each vehicle and new integrations should use /vehicles.
func (ro *Routes) FleetsRoutes(r *gin.Engine) {
	fr := fleets_handlers.NewFleetHanlder(ro.Context, ro.FleetsDB, ro.Trace, ro.Logging)
	routes := r.Group("/api/v1/fleets")
	{
		routes.POST("", fr.CreateFleets)
		routes.GET("", fr.ListFleets)
		routes.GET("/:fleetID", fr.ListFleetsById)
		routes.PATCH("/:fleetID", fr.UpdateFleets)
	}
}
