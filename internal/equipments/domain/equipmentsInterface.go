package equipments_domain

import "github.com/gin-gonic/gin"

type IEquipments interface {
	AssignEquipmentFleet(c *gin.Context)
	CreateEquipment(c *gin.Context)
	ListEquipments(c *gin.Context)
	ListEquipmentsFleets(c *gin.Context)
	ListEquipmentsById(c *gin.Context)
	ListEquipmentStatusHistory(c *gin.Context)
	RemoveEquipmentFleet(c *gin.Context)
	UpdateEquipments(c *gin.Context)
	UpdateEquipmentStatus(c *gin.Context)
}
