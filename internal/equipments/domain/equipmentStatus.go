package equipments_domain

import (
	"time"

	"github.com/google/uuid"
)

type EquipmentStatusHistory struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`

	EquipmentID uuid.UUID `json:"equipmentId" gorm:"type:uuid;not null;index"`

	FromStatus EquipmentStatus `json:"fromStatus" gorm:"size:30;not null"`
	ToStatus   EquipmentStatus `json:"toStatus" gorm:"size:30;not null"`

	Reason string `json:"reason" gorm:"size:250;not null"`

	ChangedBy *uuid.UUID `json:"changedBy,omitempty" gorm:"type:uuid"`
	ChangedAt time.Time  `json:"changedAt" gorm:"not null"`

	Equipment Equipment `json:"-" gorm:"foreignKey:EquipmentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (EquipmentStatusHistory) TableName() string {
	return "equipment_status_history"
}

var allowedStatusTransitions = map[EquipmentStatus]map[EquipmentStatus]struct{}{
	StatusAvailable: {
		StatusReserved:    {},
		StatusMaintenance: {},
		StatusInactive:    {},
	},
	StatusReserved: {
		StatusAvailable:   {},
		StatusInTransit:   {},
		StatusWorking:     {},
		StatusMaintenance: {},
	},
	StatusInTransit: {
		StatusAvailable:   {},
		StatusWorking:     {},
		StatusMaintenance: {},
	},
	StatusWorking: {
		StatusAvailable:   {},
		StatusMaintenance: {},
	},
	StatusMaintenance: {
		StatusAvailable: {},
		StatusInactive:  {},
	},
	StatusInactive: {
		StatusAvailable: {},
		StatusRetired:   {},
	},
	StatusRetired: {},
}

func (s EquipmentStatus) IsValid() bool {
	switch s {
	case StatusAvailable,
		StatusReserved,
		StatusInTransit,
		StatusWorking,
		StatusMaintenance,
		StatusInactive,
		StatusRetired:
		return true
	default:
		return false
	}
}

func (s EquipmentStatus) CanTransitionTo(next EquipmentStatus) bool {
	allowed, exists := allowedStatusTransitions[s]
	if !exists {
		return false
	}

	_, valid := allowed[next]

	return valid
}
