package equipments_db_postgres

import (
	"net/http"
	"strings"

	equipments_domain "github.com/tariq-ventura/fleet-service/internal/equipments/domain"
	"github.com/tariq-ventura/fleet-service/internal/interfaces"
)

func (pc *PostgresClient) ListEquipments(page, pageSize int, etype, status, brand, search string) ([]equipments_domain.Equipment, *interfaces.Error, int64) {
	operationSpan, spanCtx := pc.trace.StartSpan(pc.ctx, "equipments.database.postgres", map[string]any{
		"db.name":           "equipments",
		"db.operation":      "list",
		"db.type":           "postgresql",
		"equipments.type":   etype,
		"equipments.status": status,
		"equipments.brand":  brand,
		"request.page":      page,
		"rquest.pageSize":   pageSize,
	})
	defer operationSpan.End()

	if pageSize > 100 {
		pageSize = 100
	}

	query := pc.client.WithContext(spanCtx).Model(&equipments_domain.Equipment{})

	if etype != "" {
		query = query.Where(
			"type = ?",
			strings.ToUpper(strings.TrimSpace(etype)),
		)
	}

	if status != "" {
		query = query.Where(
			"status = ?",
			strings.ToUpper(strings.TrimSpace(status)),
		)
	}

	if brand != "" {
		query = query.Where(
			"brand = ?",
			strings.ToUpper(strings.TrimSpace(brand)),
		)
	}

	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where(
			`code ILIKE ? 
			OR brand ILIKE ? 
			OR model ILIKE ? 
			OR serial_number ILIKE ?`,
			pattern,
			pattern,
			pattern,
			pattern,
		)
	}

	var total int64

	if err := query.Count(&total).Error; err != nil {
		pc.logging.LogError("database_error", map[string]any{"error": err.Error()})
		return nil, &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo contar la maquinaria",
			StatusCode: http.StatusInternalServerError,
		}, 0
	}

	var equipment []equipments_domain.Equipment
	offset := (page - 1) * pageSize

	result := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&equipment)

	if result.Error != nil {
		pc.logging.LogError("database_error", map[string]any{"error": result.Error.Error()})
		return nil, &interfaces.Error{
			Error:      "database_error",
			Message:    "No se pudo contar la maquinaria",
			StatusCode: http.StatusInternalServerError,
		}, 0
	}

	return equipment, nil, total
}
