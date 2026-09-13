package geofences_dto

type CreateGeofenceRequest struct {
	GeofenceID       string  `json:"geofenceId" binding:"required,max=100"`
	Name             string  `json:"name" binding:"required,max=200"`
	Group            string  `json:"group" binding:"omitempty,max=150"`
	AdditionalMargin float64 `json:"additionalMargin" binding:"gte=0"`
	Latitude         int64   `json:"latitude" binding:"required"`
	Longitude        int64   `json:"longitude" binding:"required"`
}
