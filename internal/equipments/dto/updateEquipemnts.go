package equipments_dto

type UpdateEquipmentRequest struct {
	FleetID *string `json:"fleetId,omitempty"`

	Type *string `json:"type,omitempty" binding:"omitempty,oneof=EXCAVATOR BACKHOE BULLDOZER LOADER CRANE TRUCK"`

	Brand        *string `json:"brand,omitempty" binding:"omitempty,min=2,max=100"`
	Model        *string `json:"model,omitempty" binding:"omitempty,min=1,max=100"`
	SerialNumber *string `json:"serialNumber,omitempty" binding:"omitempty,min=3,max=100"`
	Year         *int    `json:"year,omitempty" binding:"omitempty,gte=1950,lte=2100"`

	CapacityTons *float64 `json:"capacityTons,omitempty" binding:"omitempty,gt=0"`

	Location *UpdateLocationRequest `json:"location,omitempty"`
}
