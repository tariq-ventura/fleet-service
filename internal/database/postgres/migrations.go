package database_postgres

import (
	"context"

	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	geofences_domain "github.com/tariq-ventura/fleet-service/internal/geofences/domain"
	maintenance_domain "github.com/tariq-ventura/fleet-service/internal/maintenance/domain"
	tasks_domain "github.com/tariq-ventura/fleet-service/internal/tasks/domain"
)

func (pc *PostgresClient) MigrateDatabase(ctx context.Context) error {
	err := pc.Client.AutoMigrate(
		equipments_domain.Equipment{},
		&equipments_domain.EquipmentStatusHistory{},
		&fleets_domain.Fleet{},
		&geofences_domain.Geofence{},
		&tasks_domain.Task{},
		&maintenance_domain.Maintenance{},
	)

	if err != nil {
		return err
	}

	pc.logging.LogInfo("Successfully completed migration in Postgres", nil)
	return nil
}
