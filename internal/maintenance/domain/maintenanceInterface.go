package maintenance_domain

import "github.com/gin-gonic/gin"

type IMaintenance interface {
	CreateMaintenance(c *gin.Context)
	ListMaintenance(c *gin.Context)
	ListMaintenanceById(c *gin.Context)
	UpdateMaintenance(c *gin.Context)
	RemoveMaintenance(c *gin.Context)
}
