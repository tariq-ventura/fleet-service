package maintenance_dto

import "time"

type UpdateMaintenanceRequest struct {
	Vehicle      *string    `json:"vehicle" binding:"omitempty,max=200"`
	Reference    *string    `json:"reference" binding:"omitempty,max=2000"`
	ServiceDate  *time.Time `json:"serviceDate"`
	Odometer     *float64   `json:"odometer" binding:"omitempty,gte=0"`
	ServiceTime  *string    `json:"serviceTime" binding:"omitempty,len=5"`
	HourMeter    *float64   `json:"hourMeter" binding:"omitempty,gte=0"`
	RepairReason *string    `json:"repairReason" binding:"omitempty,max=200"`
	Provider     *string    `json:"provider" binding:"omitempty,max=200"`
	Mechanic     *string    `json:"mechanic" binding:"omitempty,max=200"`
	ServiceType  *string    `json:"serviceType" binding:"omitempty,max=200"`
}
