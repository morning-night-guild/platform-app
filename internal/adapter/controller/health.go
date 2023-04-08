package controller

import (
	"context"

	"github.com/bufbuild/connect-go"
	healthv1 "github.com/morning-night-guild/platform-app/pkg/connect/health/v1"
	"github.com/morning-night-guild/platform-app/pkg/connect/health/v1/healthv1connect"
)

var _ healthv1connect.HealthServiceHandler = (*Health)(nil)

// Health.
type Health struct{}

// NewHealth ヘルスコントローラを新規作成する関数.
func NewHealth() *Health {
	return &Health{}
}

// Check ヘルスチェックするコントローラメソッド.
func (h *Health) Check(
	_ context.Context,
	_ *connect.Request[healthv1.CheckRequest],
) (*connect.Response[healthv1.CheckResponse], error) {
	return connect.NewResponse(&healthv1.CheckResponse{}), nil
}
