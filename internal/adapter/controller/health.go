package controller

import (
	"context"

	"github.com/bufbuild/connect-go"
	healthv1 "github.com/morning-night-guild/platform-app/pkg/connect/proto/health/v1"
)

// Health.
type Health struct{}

// NewHealth ヘルスコントローラを新規作成する関数.
func NewHealth() *Health {
	return &Health{}
}

// Check ヘルスチェックするコントローラメソッド.
func (h *Health) Check(
	ctx context.Context,
	req *connect.Request[healthv1.CheckRequest],
) (*connect.Response[healthv1.CheckResponse], error) {
	return connect.NewResponse(&healthv1.CheckResponse{}), nil
}
