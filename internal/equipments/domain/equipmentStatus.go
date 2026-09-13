package equipments_domain

import (
	"time"

	"github.com/google/uuid"
)

type EquipmentStatusHistory struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	EquipmentID uuid.UUID `json:"equipmentId" gorm:"type:uuid;not null;index"`
	FromStatus  string    `json:"fromStatus" gorm:"size:50;not null"`
	ToStatus    string    `json:"toStatus" gorm:"size:50;not null"`
	Reason      string    `json:"reason" gorm:"size:250;not null"`
	ChangedAt   time.Time `json:"changedAt" gorm:"not null"`

	Equipment Equipment `json:"-" gorm:"foreignKey:EquipmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (EquipmentStatusHistory) TableName() string {
	return "startrack_vehicle_status_history"
}
