package fleets_domain

import "github.com/gin-gonic/gin"

type IFleets interface {
	CreateFleets(c *gin.Context)
	ListFleets(c *gin.Context)
	ListFleetsById(c *gin.Context)
	UpdateFleets(c *gin.Context)
}
