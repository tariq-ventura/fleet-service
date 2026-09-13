package database

import (
	"context"
	"errors"

	database_postgres "github.com/tariq-ventura/fleet-service/internal/database/postgres"
	equipments_db "github.com/tariq-ventura/fleet-service/internal/equipments/db"
	fleets_db "github.com/tariq-ventura/fleet-service/internal/fleets/db"
	geofences_db "github.com/tariq-ventura/fleet-service/internal/geofences/db"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	maintenance_db "github.com/tariq-ventura/fleet-service/internal/maintenance/db"
	tasks_db "github.com/tariq-ventura/fleet-service/internal/tasks/db"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

type Database struct {
	Equipments  equipments_db.IEquipmentsDB
	Fleets      fleets_db.IFleetsDB
	Geofences   geofences_db.IGeofencesDB
	Tasks       tasks_db.ITasksDB
	Maintenance maintenance_db.IMaintenanceDB
}

type IDatabase interface {
	MigrateDatabase(ctx context.Context) error
}

var SetupDatabase = func(ctx context.Context, l logging.ILogging, t interfaces.ITrace) (*Database, IDatabase, error) {
	dbType, err := validations.RequiredEnv("DB_CONTEXT")
	if err != nil {
		return nil, nil, err
	}

	switch dbType {
	case "postgresql":
		l.LogInfo("Database selected", map[string]any{"db_type": dbType})
		db, err := database_postgres.SetupPostgres(l)

		if err != nil {
			return nil, nil, err
		}

		equipments, err := equipments_db.NewDatabase(ctx, l, t, db.Client)

		if err != nil {
			return nil, nil, err
		}

		fleets, err := fleets_db.NewDatabase(ctx, l, t, db.Client)

		if err != nil {
			return nil, nil, err
		}

		geofences, err := geofences_db.NewDatabase(ctx, l, t, db.Client)
		if err != nil {
			return nil, nil, err
		}
		tasks, err := tasks_db.NewDatabase(ctx, l, t, db.Client)
		if err != nil {
			return nil, nil, err
		}
		maintenance, err := maintenance_db.NewDatabase(ctx, l, t, db.Client)
		if err != nil {
			return nil, nil, err
		}

		return &Database{
			Equipments:  equipments,
			Fleets:      fleets,
			Geofences:   geofences,
			Tasks:       tasks,
			Maintenance: maintenance,
		}, db, err
	default:
		return nil, nil, errors.New("unsupported database backend")
	}
}
