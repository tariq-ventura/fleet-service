package equipments_handlers

import (
	"github.com/gin-gonic/gin"
	equipments_db "github.com/tariq-ventura/fleet-service/internal/equipments/db"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
)

type EquipmentHanlder struct {
	db    equipments_db.IEquipmentsDB
	trace interfaces.ITrace
	logs  logging.ILogging
}

func NewEquipmentHanlder(server *gin.Context, db equipments_db.IEquipmentsDB, trace interfaces.ITrace, logs logging.ILogging) equipments_domain.IEquipments {
	return &EquipmentHanlder{
		db:    db,
		trace: trace,
		logs:  logs,
	}
}
