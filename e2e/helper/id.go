package helper

import (
	"testing"

	"github.com/google/uuid"
)

func GenerateIDs(t *testing.T, count int) []uuid.UUID {
	t.Helper()

	ids := make([]uuid.UUID, count)

	for i := 0; i < count; i++ {
		ids[i] = uuid.New()
	}

	return ids
}
