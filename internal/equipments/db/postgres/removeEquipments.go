package equipments_db_postgres

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
)

func (pc *PostgresClient) RemoveEquipment(id uuid.UUID, ctx context.Context) *interfaces.Error {
	result := pc.client.WithContext(ctx).Delete(&equipments_domain.Equipment{}, "id = ?", id)
	if result.Error != nil {
		return databaseError("No se pudo eliminar el vehículo", result.Error, pc)
	}
	if result.RowsAffected == 0 {
		return &interfaces.Error{Error: "vehicle_not_found", Message: "El vehículo no existe", StatusCode: http.StatusNotFound}
	}

	return nil
}
