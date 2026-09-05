package database

import (
	"context"
	"errors"

	database_postgres "github.com/tariq-ventura/fleet-service/internal/database/postgres"
	equipments_db "github.com/tariq-ventura/fleet-service/internal/equipments/db"
	fleets_db "github.com/tariq-ventura/fleet-service/internal/fleets/db"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	"github.com/tariq-ventura/fleet-service/internal/validations"
)

type Database struct {
	Equipments equipments_db.IEquipmentsDB
	Fleets     fleets_db.IFleetsDB
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

		return &Database{
			Equipments: equipments,
			Fleets:     fleets,
		}, db, err
	default:
		return nil, nil, errors.New("unsupported database backend")
	}
}
