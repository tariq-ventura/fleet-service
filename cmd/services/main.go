package main

import (
	"context"
	"log"

	"github.com/tariq-ventura/fleet-service/internal/database"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	"github.com/tariq-ventura/fleet-service/internal/router"
	"github.com/tariq-ventura/fleet-service/internal/trace"
)

func initTrace(ctx context.Context, l logging.ILogging) (interfaces.ITrace, error) {
	l.LogInfo("Initializing tracing system", nil)

	trace, err := trace.NewTrace(ctx)

	if err != nil {
		l.LogError("Tracing initialization error", map[string]any{"error": err.Error()})
		return nil, err
	}

	l.LogInfo("Tracing system initialized successfully", nil)
	return trace, nil
}

func runApp(ctx context.Context, l logging.ILogging, t interfaces.ITrace) error {
	r := &router.Routes{}
	r.Logging = l
	r.Trace = t

	database, db, _ := database.SetupDatabase(ctx, l, t)

	err := db.MigrateDatabase(ctx)

	if err != nil {
		return err
	}

	r.EquipmentsDB = database.Equipments
	r.FleetsDB = database.Fleets

	r.Routes = r.SetupRouter()
	r.Run()

	return nil
}

func main() {
	ctx := context.Background()
	logging, err := logging.NewLogging(ctx)

	if err != nil {
		log.Fatalf("Logging setup error: %v", err)
		return
	}

	logging.LogInfo("Monitoring initialized successfully", nil)

	trace, err := initTrace(ctx, logging)

	if err != nil {
		logging.LogError("Trace setup error", map[string]any{"error": err.Error()})
		return
	}

	if err := runApp(ctx, logging, trace); err != nil {
		logging.LogError("Router setup error: ", map[string]any{"error": err.Error()})
		return
	}
}
