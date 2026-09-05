package fleets_db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	fleets_db_postgres "github.com/tariq-ventura/fleet-service/internal/fleets/db/postgres"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	"github.com/tariq-ventura/fleet-service/internal/validations"
	"gorm.io/gorm"
)

type IFleetsDB interface {
	CreateFleets(fleet fleets_domain.Fleet) *interfaces.Error
	ListFleets(page, pageSize int, search string) ([]fleets_domain.Fleet, *interfaces.Error, int64)
	ListFleetsById(id uuid.UUID) (*fleets_domain.Fleet, *interfaces.Error)
	UpdateFleets(id uuid.UUID, updates map[string]any) (*fleets_domain.Fleet, *interfaces.Error)
}

var NewDatabase = func(ctx context.Context, l logging.ILogging, t interfaces.ITrace, client *gorm.DB) (IFleetsDB, error) {
	dbType, err := validations.RequiredEnv("DB_CONTEXT")
	if err != nil {
		return nil, err
	}

	switch dbType {
	case "postgresql":
		l.LogInfo("Database selected", map[string]any{"db_type": dbType})
		return fleets_db_postgres.SetupPostgres(ctx, l, t, client)
	default:
		return nil, errors.New("unsupported database backend")
	}
}
