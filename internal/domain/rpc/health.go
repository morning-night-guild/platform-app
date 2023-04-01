package rpc

import "context"

//go:generate mockgen -source health.go -destination health_mock.go -package rpc

type Health interface {
	Check(context.Context) error
}
