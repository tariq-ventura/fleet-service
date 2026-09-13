package router

import (
	"github.com/gin-gonic/gin"
	tasks_handlers "github.com/tariq-ventura/fleet-service/internal/tasks/handlers"
)

func (ro *Routes) TasksRoutes(r *gin.Engine) {
	h := tasks_handlers.NewTaskHandler(ro.Context, ro.TasksDB, ro.Trace, ro.Logging)
	routes := r.Group("/api/v1/tasks")
	{
		routes.POST("", h.CreateTask)
		routes.GET("", h.ListTasks)
		routes.GET("/:id", h.ListTaskById)
		routes.PATCH("/:id", h.UpdateTask)
		routes.DELETE("/:id", h.RemoveTask)
	}
}
