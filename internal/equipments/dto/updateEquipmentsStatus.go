package equipments_dto

type UpdateEquipmentStatusRequest struct {
	Status string `json:"status" binding:"required,min=2,max=50"`
	Reason string `json:"reason" binding:"required,min=2,max=250"`
}
