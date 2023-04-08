package usecase

import (
	"context"
)

//go:generate mockgen -source api_health.go -destination api_health_mock.go -package usecase

// APIHealth.
type APIHealth interface {
	Check(context.Context, APIHealthCheckInput) (APIHealthCheckOutput, error)
}

// APIHealthCheckInput.
type APIHealthCheckInput struct{}

// APIHealthCheckOutput.
type APIHealthCheckOutput struct{}
