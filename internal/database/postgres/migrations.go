package database_postgres

import (
	"context"

	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
)

func (pc *PostgresClient) MigrateDatabase(ctx context.Context) error {
	err := pc.Client.AutoMigrate(
		equipments_domain.Equipment{},
		&equipments_domain.EquipmentStatusHistory{},
		&fleets_domain.Fleet{},
	)

	if err != nil {
		return err
	}

	pc.logging.LogInfo("Successfully completed migration in Postgres", nil)
	return nil
}
