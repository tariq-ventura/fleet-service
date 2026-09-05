package fleets_db_postgres

import (
	"net/http"

	fleets_domain "github.com/tariq-ventura/fleet-service/internal/fleets/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
)

func (pc *PostgresClient) ListFleets(page, pageSize int, search string) ([]fleets_domain.Fleet, *interfaces.Error, int64) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":         "fleets",
		"db.operation":    "list",
		"db.type":         "postgresql",
		"request.page":    page,
		"rquest.pageSize": pageSize,
	})
	defer operationSpan.End()

	if pageSize > 100 {
		pageSize = 100
	}

	query := pc.client.WithContext(spanCtx).Model(&fleets_domain.Fleet{})

	if search != "" {
		pattern := "%" + search + "%"

		query = query.Where(
			"code ILIKE ? OR name ILIKE ?",
			pattern,
			pattern,
		)
	}

	var total int64

	if err := query.Count(&total).Error; err != nil {
		pc.logging.LogError("database_error", map[string]any{"error": err.Error()})
		return nil, &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo contar las flotas",
			StatusCode: http.StatusInternalServerError,
		}, 0
	}

	var fleets []fleets_domain.Fleet
	offset := (page - 1) * pageSize

	result := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&fleets)

	if result.Error != nil {
		pc.logging.LogError("database_error", map[string]any{"error": result.Error.Error()})
		return nil, &interfaces.Error{
			Error:      "fleet_database_error",
			Message:    "No se pudo contar las flotas",
			StatusCode: http.StatusInternalServerError,
		}, 0
	}

	return fleets, nil, total
}
