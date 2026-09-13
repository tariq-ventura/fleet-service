package geofences_domain

import "github.com/gin-gonic/gin"

type IGeofences interface {
	CreateGeofence(c *gin.Context)
	ListGeofences(c *gin.Context)
	ListGeofenceById(c *gin.Context)
	UpdateGeofence(c *gin.Context)
	RemoveGeofence(c *gin.Context)
}
