package geofences_db_postgres

import (
	"context"

	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	"gorm.io/gorm"
)

type PostgresClient struct {
	client  *gorm.DB
	logging logging.ILogging
	trace   interfaces.ITrace
}

func SetupPostgres(_ context.Context, logs logging.ILogging, trace interfaces.ITrace, client *gorm.DB) (*PostgresClient, error) {
	return &PostgresClient{client: client, logging: logs, trace: trace}, nil
}
