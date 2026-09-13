package equipments_domain

import "github.com/gin-gonic/gin"

type IEquipments interface {
	CreateEquipment(c *gin.Context)
	ListEquipments(c *gin.Context)
	ListEquipmentsById(c *gin.Context)
	UpdateEquipments(c *gin.Context)
	RemoveEquipment(c *gin.Context)
	UpdateEquipmentStatus(c *gin.Context)
	ListEquipmentStatusHistory(c *gin.Context)
}
