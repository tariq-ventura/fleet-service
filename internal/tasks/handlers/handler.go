package tasks_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"github.com/tariq-ventura/fleet-service/internal/logging"
	tasks_db "github.com/tariq-ventura/fleet-service/internal/tasks/db"
	tasks_domain "github.com/tariq-ventura/fleet-service/internal/tasks/domain"
)

type TaskHandler struct {
	db    tasks_db.ITasksDB
	trace interfaces.ITrace
	logs  logging.ILogging
}

func NewTaskHandler(_ *gin.Context, db tasks_db.ITasksDB, trace interfaces.ITrace, logs logging.ILogging) tasks_domain.ITasks {
	return &TaskHandler{db: db, trace: trace, logs: logs}
}
