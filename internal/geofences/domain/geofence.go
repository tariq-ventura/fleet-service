package geofences_domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Geofence struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	GeofenceID       string         `json:"geofenceId" gorm:"size:100;not null;uniqueIndex"`
	Name             string         `json:"name" gorm:"size:200;not null;index"`
	Group            string         `json:"group" gorm:"size:150"`
	AdditionalMargin float64        `json:"additionalMargin" gorm:"not null;default:0"`
	Latitude         int64          `json:"latitude" gorm:"not null"`
	Longitude        int64          `json:"longitude" gorm:"not null"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Geofence) TableName() string { return "startrack_geofences" }

func (g *Geofence) BeforeCreate(_ *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}
