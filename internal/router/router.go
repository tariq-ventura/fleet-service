package router

import (
	"github.com/gin-gonic/gin"
	equipments_db "github.com/tariq-ventura/fleet-service/internal/equipments/db"
	fleets_db "github.com/tariq-ventura/fleet-service/internal/fleets/db"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
)

type Routes struct {
	Routes       *gin.Engine
	Context      *gin.Context
	Logging      logging.ILogging
	Trace        interfaces.ITrace
	EquipmentsDB equipments_db.IEquipmentsDB
	FleetsDB     fleets_db.IFleetsDB
}

func (r *Routes) SetupRouter() *gin.Engine {
	r.Routes = gin.Default()

	r.SetupCors()

	r.HealthCheckRoutes()
	r.EquipmentsRoutes(r.Routes)
	r.FleetsRoutes(r.Routes)
	return r.Routes
}

func (r *Routes) Run() {
	r.Routes.Run(":3000")
}
