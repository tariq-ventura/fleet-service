package equipments_db_postgres

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (pc *PostgresClient) UpdateEquipmentStatus(id uuid.UUID, status, reason string, ctx context.Context) (*equipments_domain.EquipmentStatusHistory, *equipments_domain.Equipment, *interfaces.Error) {
	var vehicle equipments_domain.Equipment
	var history equipments_domain.EquipmentStatusHistory

	err := pc.client.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&vehicle, "id = ?", id).Error; err != nil {
			return err
		}

		next := strings.TrimSpace(status)
		if strings.EqualFold(vehicle.Status, next) {
			return errors.New("vehicle already has the requested status")
		}

		previous := vehicle.Status
		if err := tx.Model(&vehicle).Update("status", next).Error; err != nil {
			return err
		}

		history = equipments_domain.EquipmentStatusHistory{
			ID: uuid.New(), EquipmentID: id, FromStatus: previous, ToStatus: next,
			Reason: strings.TrimSpace(reason), ChangedAt: time.Now().UTC(),
		}
		return tx.Create(&history).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, &interfaces.Error{Error: "vehicle_not_found", Message: "El vehículo no existe", StatusCode: http.StatusNotFound}
	}
	if err != nil {
		return nil, nil, &interfaces.Error{Error: "invalid_status_change", Message: err.Error(), StatusCode: http.StatusConflict}
	}

	return &history, &vehicle, nil
}
