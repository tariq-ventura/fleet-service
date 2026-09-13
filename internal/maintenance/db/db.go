package maintenance_db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	maintenance_db_postgres "github.com/tariq-ventura/fleet-service/internal/maintenance/db/postgres"
	maintenance_domain "github.com/tariq-ventura/fleet-service/internal/maintenance/domain"
	"github.com/tariq-ventura/fleet-service/internal/validations"
	"gorm.io/gorm"
)

type IMaintenanceDB interface {
	CreateMaintenance(data *maintenance_domain.Maintenance, ctx context.Context) *interfaces.Error
	ListMaintenance(page, pageSize int, vehicle, search string, ctx context.Context) ([]maintenance_domain.Maintenance, *interfaces.Error, int64)
	ListMaintenanceById(id uuid.UUID, ctx context.Context) (*maintenance_domain.Maintenance, *interfaces.Error)
	UpdateMaintenance(id uuid.UUID, updates map[string]any, ctx context.Context) (*maintenance_domain.Maintenance, *interfaces.Error)
	RemoveMaintenance(id uuid.UUID, ctx context.Context) *interfaces.Error
}

func NewDatabase(ctx context.Context, l logging.ILogging, t interfaces.ITrace, client *gorm.DB) (IMaintenanceDB, error) {
	if dbType, err := validations.RequiredEnv("DB_CONTEXT"); err != nil || dbType != "postgresql" {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("unsupported database backend")
	}
	return maintenance_db_postgres.SetupPostgres(ctx, l, t, client)
}
