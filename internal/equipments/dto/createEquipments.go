package equipments_dto

type CreateEquipmentRequest struct {
	Code string `json:"code" binding:"required,min=3,max=50"`

	FleetID *string `json:"fleetId,omitempty"`

	Type string `json:"type" binding:"required,oneof=EXCAVATOR BACKHOE BULLDOZER LOADER CRANE TRUCK"`

	Brand        string `json:"brand" binding:"required,min=2,max=100"`
	Model        string `json:"model" binding:"required,min=1,max=100"`
	SerialNumber string `json:"serialNumber" binding:"required,min=3,max=100"`
	Year         int    `json:"year" binding:"required,gte=1950,lte=2100"`

	CapacityTons float64 `json:"capacityTons" binding:"required,gt=0"`

	Location CreateLocationRequest `json:"location" binding:"required"`

	EngineHours         float64 `json:"engineHours" binding:"gte=0"`
	MaintenanceInterval float64 `json:"maintenanceIntervalHours" binding:"required,gt=0"`
	FuelPercent         float64 `json:"fuelPercent" binding:"gte=0,lte=100"`
}
