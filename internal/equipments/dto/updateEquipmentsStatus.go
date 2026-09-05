package equipments_dto

type UpdateEquipmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=AVAILABLE RESERVED IN_TRANSIT WORKING MAINTENANCE INACTIVE RETIRED"`
	Reason string `json:"reason" binding:"required,min=3,max=250"`
}
