package tasks_domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey"`
	TaskID        string         `json:"taskId" gorm:"size:100;not null;uniqueIndex"`
	Title         string         `json:"title" gorm:"size:200;not null"`
	Description   string         `json:"description" gorm:"type:text"`
	Type          string         `json:"type" gorm:"size:100;not null;index"`
	Status        string         `json:"status" gorm:"size:100;not null;index"`
	ScheduledDate time.Time      `json:"scheduledDate" gorm:"not null;index"`
	Origin        string         `json:"origin" gorm:"size:250;not null"`
	Destination   string         `json:"destination" gorm:"size:250;not null"`
	Latitude      int64          `json:"latitude" gorm:"not null"`
	Longitude     int64          `json:"longitude" gorm:"not null"`
	Assignee      string         `json:"assignee" gorm:"size:200;not null"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Task) TableName() string { return "startrack_tasks" }

func (t *Task) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
