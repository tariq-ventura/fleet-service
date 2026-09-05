package equipments_db_postgres

import (
	"context"

	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
)

func (pc *PostgresClient) MigrateEquipments(ctx context.Context) error {
	err := pc.client.AutoMigrate(equipments_domain.Equipment{}, &equipments_domain.EquipmentStatusHistory{})

	if err != nil {
		return err
	}

	pc.logging.LogInfo("Successfully completed migration in Postgres", nil)
	return nil
}
