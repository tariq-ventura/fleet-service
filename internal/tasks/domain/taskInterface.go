package tasks_domain

import "github.com/gin-gonic/gin"

type ITasks interface {
	CreateTask(c *gin.Context)
	ListTasks(c *gin.Context)
	ListTaskById(c *gin.Context)
	UpdateTask(c *gin.Context)
	RemoveTask(c *gin.Context)
}
