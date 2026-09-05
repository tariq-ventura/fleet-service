package equipments_db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	equipments_db_postgres "github.com/tariq-ventura/fleet-service/internal/equipments/db/postgres"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	"github.com/tariq-ventura/fleet-service/internal/validations"
	"gorm.io/gorm"
)

type IEquipmentsDB interface {
	AssignEquipmentFleet(fleetId, equipmentId uuid.UUID) (*equipments_domain.Equipment, *interfaces.Error)
	CreateEquipment(data equipments_domain.Equipment) *interfaces.Error
	ListEquipments(page, pageSize int, etype, status, brand, search string) ([]equipments_domain.Equipment, *interfaces.Error, int64)
	ListEquipmentsById(id uuid.UUID) (*equipments_domain.Equipment, *interfaces.Error)
	ListEquipmentStatusHistory(id uuid.UUID) ([]equipments_domain.EquipmentStatusHistory, *interfaces.Error)
	ListEquipmentsFleets(fleeId uuid.UUID) ([]equipments_domain.Equipment, *fleets_domain.Fleet, *interfaces.Error)
	RemoveEquipmentFleet(fleetId, equipmentId uuid.UUID) *interfaces.Error
	UpdateEquipments(id uuid.UUID, updates map[string]any) (*equipments_domain.Equipment, *interfaces.Error)
	UpdateEquipmentStatus(equipmentId uuid.UUID, newStatus equipments_domain.EquipmentStatus) (*equipments_domain.EquipmentStatusHistory, *equipments_domain.Equipment, *interfaces.Error)
}

var NewDatabase = func(ctx context.Context, l logging.ILogging, t interfaces.ITrace, client *gorm.DB) (IEquipmentsDB, error) {
	dbType, err := validations.RequiredEnv("DB_CONTEXT")
	if err != nil {
		return nil, err
	}

	switch dbType {
	case "postgresql":
		l.LogInfo("Database selected", map[string]any{"db_type": dbType})
		return equipments_db_postgres.SetupPostgres(ctx, l, t, client)
	default:
		return nil, errors.New("unsupported database backend")
	}
}
