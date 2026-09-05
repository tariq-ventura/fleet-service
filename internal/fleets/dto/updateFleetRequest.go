package fleets_dto

type UpdateFleetRequest struct {
	Code *string `json:"code,omitempty" binding:"omitempty,min=3,max=50"`
	Name *string `json:"name,omitempty" binding:"omitempty,min=3,max=150"`
}
