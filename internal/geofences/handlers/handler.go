package geofences_handlers

import (
	"github.com/gin-gonic/gin"
	geofences_db "github.com/tariq-ventura/fleet-service/internal/geofences/db"
	geofences_domain "github.com/tariq-ventura/fleet-service/internal/geofences/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
)

type GeofenceHandler struct {
	db    geofences_db.IGeofencesDB
	trace interfaces.ITrace
	logs  logging.ILogging
}

func NewGeofenceHandler(_ *gin.Context, db geofences_db.IGeofencesDB, trace interfaces.ITrace, logs logging.ILogging) geofences_domain.IGeofences {
	return &GeofenceHandler{db: db, trace: trace, logs: logs}
}
