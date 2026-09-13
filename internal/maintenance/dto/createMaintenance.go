package maintenance_dto

import "time"

type CreateMaintenanceRequest struct {
	Vehicle      string    `json:"vehicle" binding:"required,max=200"`
	Reference    string    `json:"reference" binding:"required,max=2000"`
	ServiceDate  time.Time `json:"serviceDate" binding:"required"`
	Odometer     float64   `json:"odometer" binding:"gte=0"`
	ServiceTime  string    `json:"serviceTime" binding:"required,len=5"`
	HourMeter    float64   `json:"hourMeter" binding:"gte=0"`
	RepairReason string    `json:"repairReason" binding:"required,max=200"`
	Provider     string    `json:"provider" binding:"omitempty,max=200"`
	Mechanic     string    `json:"mechanic" binding:"omitempty,max=200"`
	ServiceType  string    `json:"serviceType" binding:"omitempty,max=200"`
}
