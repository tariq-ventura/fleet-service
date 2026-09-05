package fleets_domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Fleet struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`

	Code        string `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Name        string `json:"name" gorm:"size:150;not null"`
	Description string `json:"description" gorm:"size:500"`
	Branch      string `json:"branch" gorm:"size:150;not null"`
	Active      bool   `json:"active" gorm:"not null;default:true"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Fleet) TableName() string {
	return "fleets"
}

func (f *Fleet) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}

	return nil
}
