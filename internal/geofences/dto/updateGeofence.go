package geofences_dto

type UpdateGeofenceRequest struct {
	GeofenceID       *string  `json:"geofenceId" binding:"omitempty,max=100"`
	Name             *string  `json:"name" binding:"omitempty,max=200"`
	Group            *string  `json:"group" binding:"omitempty,max=150"`
	AdditionalMargin *float64 `json:"additionalMargin" binding:"omitempty,gte=0"`
	Latitude         *int64   `json:"latitude"`
	Longitude        *int64   `json:"longitude"`
}
