package maintenance_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	maintenance_db "github.com/tariq-ventura/fleet-service/internal/maintenance/db"
	maintenance_domain "github.com/tariq-ventura/fleet-service/internal/maintenance/domain"
)

type MaintenanceHandler struct {
	db    maintenance_db.IMaintenanceDB
	trace interfaces.ITrace
	logs  logging.ILogging
}

func NewMaintenanceHandler(_ *gin.Context, db maintenance_db.IMaintenanceDB, trace interfaces.ITrace, logs logging.ILogging) maintenance_domain.IMaintenance {
	return &MaintenanceHandler{db: db, trace: trace, logs: logs}
}
