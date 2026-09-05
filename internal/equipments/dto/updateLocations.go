package equipments_dto

type UpdateLocationRequest struct {
	Name      *string  `json:"name,omitempty" binding:"omitempty,min=2,max=150"`
	Latitude  *float64 `json:"latitude,omitempty" binding:"omitempty,gte=-90,lte=90"`
	Longitude *float64 `json:"longitude,omitempty" binding:"omitempty,gte=-180,lte=180"`
}
