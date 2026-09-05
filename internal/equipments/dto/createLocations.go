package equipments_dto

type CreateLocationRequest struct {
	Name      string  `json:"name" binding:"required,min=2,max=150"`
	Latitude  float64 `json:"latitude" binding:"gte=-90,lte=90"`
	Longitude float64 `json:"longitude" binding:"gte=-180,lte=180"`
}
