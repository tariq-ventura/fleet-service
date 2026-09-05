package fleets_dto

type CreateFleetRequest struct {
	Code string `json:"code" binding:"required,min=3,max=50"`
	Name string `json:"name" binding:"required,min=3,max=150"`
}
