package equipments_dto

type UpdateEquipmentRequest struct {
	Description *string `json:"description" binding:"omitempty,min=2,max=200"`
	Type        *string `json:"type" binding:"omitempty,min=2,max=100"`
	Year        *int    `json:"year" binding:"omitempty,gte=1900,lte=2100"`
	Color       *string `json:"color" binding:"omitempty,max=80"`
	Brand       *string `json:"brand" binding:"omitempty,max=100"`
	Model       *string `json:"model" binding:"omitempty,max=100"`
	Group       *string `json:"group" binding:"omitempty,max=150"`
	Tags        *string `json:"tags" binding:"omitempty,max=500"`
	Driver      *string `json:"driver" binding:"omitempty,max=200"`
	RemoteID    *string `json:"remoteId" binding:"omitempty,max=100"`
}
