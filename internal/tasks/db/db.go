package tasks_db

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	tasks_db_postgres "github.com/tariq-ventura/fleet-service/internal/tasks/db/postgres"
	tasks_domain "github.com/tariq-ventura/fleet-service/internal/tasks/domain"
	"github.com/tariq-ventura/fleet-service/internal/validations"
	"gorm.io/gorm"
)

type ITasksDB interface {
	CreateTask(data *tasks_domain.Task, ctx context.Context) *interfaces.Error
	ListTasks(page, pageSize int, status, taskType, search string, ctx context.Context) ([]tasks_domain.Task, *interfaces.Error, int64)
	ListTaskById(id uuid.UUID, ctx context.Context) (*tasks_domain.Task, *interfaces.Error)
	UpdateTask(id uuid.UUID, updates map[string]any, ctx context.Context) (*tasks_domain.Task, *interfaces.Error)
	RemoveTask(id uuid.UUID, ctx context.Context) *interfaces.Error
}

func NewDatabase(ctx context.Context, l logging.ILogging, t interfaces.ITrace, client *gorm.DB) (ITasksDB, error) {
	if dbType, err := validations.RequiredEnv("DB_CONTEXT"); err != nil || dbType != "postgresql" {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("unsupported database backend")
	}
	return tasks_db_postgres.SetupPostgres(ctx, l, t, client)
}
