package fleets_handlers

import (
	"github.com/gin-gonic/gin"
	fleets_db "github.com/tariq-ventura/fleet-service/internal/fleets/db"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
)

type FleetHanlder struct {
	db    fleets_db.IFleetsDB
	trace interfaces.ITrace
	logs  logging.ILogging
}

func NewFleetHanlder(server *gin.Context, db fleets_db.IFleetsDB, trace interfaces.ITrace, logs logging.ILogging) fleets_domain.IFleets {
	return &FleetHanlder{
		db:    db,
		trace: trace,
		logs:  logs,
	}
}
