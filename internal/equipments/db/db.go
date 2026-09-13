package equipments_db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	equipments_db_postgres "github.com/tariq-ventura/fleet-service/internal/equipments/db/postgres"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	"github.com/tariq-ventura/fleet-service/internal/validations"
	"gorm.io/gorm"
)

type IEquipmentsDB interface {
	CreateEquipment(data *equipments_domain.Equipment, ctx context.Context) *interfaces.Error
	ListEquipments(page, pageSize int, equipmentType, status, brand, search string, ctx context.Context) ([]equipments_domain.Equipment, *interfaces.Error, int64)
	ListEquipmentsById(id uuid.UUID, ctx context.Context) (*equipments_domain.Equipment, *interfaces.Error)
	UpdateEquipments(id uuid.UUID, updates map[string]any, ctx context.Context) (*equipments_domain.Equipment, *interfaces.Error)
	RemoveEquipment(id uuid.UUID, ctx context.Context) *interfaces.Error
	UpdateEquipmentStatus(id uuid.UUID, status, reason string, ctx context.Context) (*equipments_domain.EquipmentStatusHistory, *equipments_domain.Equipment, *interfaces.Error)
	ListEquipmentStatusHistory(id uuid.UUID, ctx context.Context) ([]equipments_domain.EquipmentStatusHistory, *interfaces.Error)
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
