package interactor

import (
	"context"

	"github.com/morning-night-guild/platform-app/internal/domain/cache"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
)

var _ port.APIAuthGenerateCode = (*APIAuthGenerateCode)(nil)

type APIAuthGenerateCode struct {
	secret    auth.Secret
	codeCache cache.Cache[model.Code]
}

func NewAPIAuthGenerateCode(
	secret auth.Secret,
	codeCache cache.Cache[model.Code],
) *APIAuthGenerateCode {
	return &APIAuthGenerateCode{
		secret:    secret,
		codeCache: codeCache,
	}
}

func (aas *APIAuthGenerateCode) Execute(
	ctx context.Context,
	input port.APIAuthGenerateCodeInput,
) (port.APIAuthGenerateCodeOutput, error) {
	sid := input.SessionToken.GetID(aas.secret)

	code := model.GenerateCode(sid)

	if err := aas.codeCache.Set(ctx, sid.String(), code, model.DefaultCodeExpiresIn); err != nil {
		return port.APIAuthGenerateCodeOutput{}, err
	}

	return port.APIAuthGenerateCodeOutput{
		Code: code,
	}, nil
}
