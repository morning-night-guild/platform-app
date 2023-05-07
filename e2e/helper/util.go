package helper

import "testing"

func ToStringPointer(t *testing.T, value string) *string {
	t.Helper()

	return &value
}

func ToIntPointer(t *testing.T, value int) *int {
	t.Helper()

	return &value
}
