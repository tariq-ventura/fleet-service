package geofences_db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	geofences_db_postgres "github.com/tariq-ventura/fleet-service/internal/geofences/db/postgres"
	geofences_domain "github.com/tariq-ventura/fleet-service/internal/geofences/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	"github.com/tariq-ventura/fleet-service/internal/validations"
	"gorm.io/gorm"
)

type IGeofencesDB interface {
	CreateGeofence(data *geofences_domain.Geofence, ctx context.Context) *interfaces.Error
	ListGeofences(page, pageSize int, search string, ctx context.Context) ([]geofences_domain.Geofence, *interfaces.Error, int64)
	ListGeofenceById(id uuid.UUID, ctx context.Context) (*geofences_domain.Geofence, *interfaces.Error)
	UpdateGeofence(id uuid.UUID, updates map[string]any, ctx context.Context) (*geofences_domain.Geofence, *interfaces.Error)
	RemoveGeofence(id uuid.UUID, ctx context.Context) *interfaces.Error
}

func NewDatabase(ctx context.Context, l logging.ILogging, t interfaces.ITrace, client *gorm.DB) (IGeofencesDB, error) {
	if dbType, err := validations.RequiredEnv("DB_CONTEXT"); err != nil || dbType != "postgresql" {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("unsupported database backend")
	}
	return geofences_db_postgres.SetupPostgres(ctx, l, t, client)
}
