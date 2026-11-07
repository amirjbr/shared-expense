package utils

import "github.com/google/uuid"

func UniqueUUIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))

	result := make([]uuid.UUID, 0, len(ids))

	for _, id := range ids {
		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = struct{}{}
		result = append(result, id)
	}

	return result
}
