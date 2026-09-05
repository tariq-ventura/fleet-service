package equipments_domain

import (
	"time"

	"github.com/google/uuid"
)

type EquipmentStatus string
type EquipmentType string

const (
	StatusAvailable   EquipmentStatus = "AVAILABLE"
	StatusReserved    EquipmentStatus = "RESERVED"
	StatusInTransit   EquipmentStatus = "IN_TRANSIT"
	StatusWorking     EquipmentStatus = "WORKING"
	StatusMaintenance EquipmentStatus = "MAINTENANCE"
	StatusInactive    EquipmentStatus = "INACTIVE"
	StatusRetired     EquipmentStatus = "RETIRED"
)

const (
	TypeExcavator EquipmentType = "EXCAVATOR"
	TypeBackhoe   EquipmentType = "BACKHOE"
	TypeBulldozer EquipmentType = "BULLDOZER"
	TypeLoader    EquipmentType = "LOADER"
	TypeCrane     EquipmentType = "CRANE"
	TypeTruck     EquipmentType = "TRUCK"
)

type Location struct {
	Name      string  `json:"name" gorm:"column:location_name;size:150;not null"`
	Latitude  float64 `json:"latitude" gorm:"column:latitude;not null"`
	Longitude float64 `json:"longitude" gorm:"column:longitude;not null"`
}

type Equipment struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`

	Code         string        `json:"code" gorm:"size:50;not null;uniqueIndex"`
	FleetID      *uuid.UUID    `json:"fleetId,omitempty" gorm:"type:uuid;index"`
	Type         EquipmentType `json:"type" gorm:"size:30;not null;index"`
	Brand        string        `json:"brand" gorm:"size:100;not null"`
	Model        string        `json:"model" gorm:"size:100;not null"`
	SerialNumber string        `json:"serialNumber" gorm:"size:100;not null;uniqueIndex"`
	Year         int           `json:"year" gorm:"not null"`

	CapacityTons float64         `json:"capacityTons" gorm:"type:numeric(8,2);not null"`
	Status       EquipmentStatus `json:"status" gorm:"size:30;not null;default:AVAILABLE;index"`

	Location Location `json:"location" gorm:"embedded"`

	EngineHours          float64 `json:"engineHours" gorm:"type:numeric(12,2);not null;default:0"`
	NextMaintenanceHours float64 `json:"nextMaintenanceHours" gorm:"type:numeric(12,2);not null"`
	FuelPercent          float64 `json:"fuelPercent" gorm:"type:numeric(5,2);not null"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Equipment) TableName() string {
	return "equipment"
}
