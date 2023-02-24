package port

import (
	"github.com/morning-night-guild/platform-app/internal/usecase"
)

// APIHealthCheckInput.
type APIHealthCheckInput struct {
	usecase.Input
}

// APIHealthCheckOutput.
type APIHealthCheckOutput struct {
	usecase.Output
}

// APIHealthCheck.
type APIHealthCheck interface {
	usecase.Usecase[APIHealthCheckInput, APIHealthCheckOutput]
}
