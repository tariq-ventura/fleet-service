package fleets_db_postgres

import (
	"context"

	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
)

func (pc *PostgresClient) MigrateFleets(ctx context.Context) error {
	err := pc.client.AutoMigrate(fleets_domain.Fleet{})

	if err != nil {
		return err
	}

	pc.logging.LogInfo("Successfully completed migration in Postgres", nil)
	return nil
}
