package maintenance_domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Maintenance struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	Vehicle      string         `json:"vehicle" gorm:"size:200;not null;index"`
	Reference    string         `json:"reference" gorm:"type:text;not null"`
	ServiceDate  time.Time      `json:"serviceDate" gorm:"not null;index"`
	Odometer     float64        `json:"odometer" gorm:"not null"`
	ServiceTime  string         `json:"serviceTime" gorm:"size:5;not null"`
	HourMeter    float64        `json:"hourMeter" gorm:"not null"`
	RepairReason string         `json:"repairReason" gorm:"size:200;not null"`
	Provider     string         `json:"provider,omitempty" gorm:"size:200"`
	Mechanic     string         `json:"mechanic,omitempty" gorm:"size:200"`
	ServiceType  string         `json:"serviceType,omitempty" gorm:"size:200"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Maintenance) TableName() string { return "startrack_maintenance" }
func (m *Maintenance) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
