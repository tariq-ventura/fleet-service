package equipments_domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Equipment represents the vehicle/asset record exposed by Startrack.
// The package name is kept for backwards compatibility with the existing API.
type Equipment struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`

	Description    string     `json:"description" gorm:"size:200;not null;uniqueIndex"`
	Status         string     `json:"status" gorm:"size:50;not null;index"`
	Type           string     `json:"type" gorm:"size:100;not null;index"`
	Year           int        `json:"year" gorm:"not null"`
	Color          string     `json:"color" gorm:"size:80"`
	Brand          string     `json:"brand" gorm:"size:100"`
	Model          string     `json:"model" gorm:"size:100"`
	Group          string     `json:"group" gorm:"size:150;index"`
	Tags           string     `json:"tags" gorm:"size:500"`
	Driver         string     `json:"driver" gorm:"size:200"`
	RemoteID       string     `json:"remoteId" gorm:"size:100;not null;uniqueIndex"`
	Location       string     `json:"location" gorm:"size:250"`
	Latitude       float64    `json:"latitude" gorm:"type:numeric(10,7)"`
	Longitude      float64    `json:"longitude" gorm:"type:numeric(10,7)"`
	LastPositionAt *time.Time `json:"lastPositionAt,omitempty" gorm:"index"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Equipment) TableName() string {
	return "startrack_vehicles"
}

func (e *Equipment) BeforeCreate(_ *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}

	return nil
}
