package validations

import equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"

func CanChangeFleet(status equipments_domain.EquipmentStatus) bool {
	switch status {
	case equipments_domain.StatusAvailable,
		equipments_domain.StatusMaintenance,
		equipments_domain.StatusInactive:
		return true
	default:
		return false
	}
}
