package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	equipments_db "github.com/tariq-ventura/fleet-service/internal/equipments/db"
	fleets_db "github.com/tariq-ventura/fleet-service/internal/fleets/db"
	geofences_db "github.com/tariq-ventura/fleet-service/internal/geofences/db"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	maintenance_db "github.com/tariq-ventura/fleet-service/internal/maintenance/db"
	tasks_db "github.com/tariq-ventura/fleet-service/internal/tasks/db"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

type Routes struct {
	Routes        *gin.Engine
	Context       *gin.Context
	Logging       logging.ILogging
	Trace         interfaces.ITrace
	EquipmentsDB  equipments_db.IEquipmentsDB
	FleetsDB      fleets_db.IFleetsDB
	GeofencesDB   geofences_db.IGeofencesDB
	TasksDB       tasks_db.ITasksDB
	MaintenanceDB maintenance_db.IMaintenanceDB
}

func (r *Routes) SetupRouter() *gin.Engine {
	r.Routes = gin.Default()

	r.SetupCors()

	r.HealthCheckRoutes()
	r.EquipmentsRoutes(r.Routes)
	r.VehiclesRoutes(r.Routes)
	r.FleetsRoutes(r.Routes)
	r.GeofencesRoutes(r.Routes)
	r.TasksRoutes(r.Routes)
	r.MaintenanceRoutes(r.Routes)
	return r.Routes
}

func (r *Routes) Run() {
	const defaultPort = "3000"
	port, err := validations.RequiredEnv("PORT")
	if err != nil {
		port = defaultPort
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	if err := r.Routes.Run(port); err != nil {
		r.Logging.LogError("No se pudo iniciar el servidor HTTP", map[string]any{"address": port, "error": err.Error()})
	}
}
