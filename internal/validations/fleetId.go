package validations

import (
	"strings"

	"github.com/google/uuid"
)

func ValidateFleetId(fleetID *string) *uuid.UUID {
	if fleetID != nil && strings.TrimSpace(*fleetID) != "" {
		parsedFleetID, err := uuid.Parse(*fleetID)
		if err != nil {
			return nil
		}
		return &parsedFleetID
	}

	return nil
}
