package usecase

import (
	"context"
	"crypto/rsa"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
)

//go:generate mockgen -source api_auth.go -destination api_auth_mock.go -package usecase

// APIAuth.
type APIAuth interface {
	SignUp(context.Context, APIAuthSignUpInput) (APIAuthSignUpOutput, error)
	SignIn(context.Context, APIAuthSignInInput) (APIAuthSignInOutput, error)
	SignOut(context.Context, APIAuthSignOutInput) (APIAuthSignOutOutput, error)
	Verify(context.Context, APIAuthVerifyInput) (APIAuthVerifyOutput, error)
	GenerateCode(context.Context, APIAuthGenerateCodeInput) (APIAuthGenerateCodeOutput, error)
	Refresh(context.Context, APIAuthRefreshInput) (APIAuthRefreshOutput, error)
}

type APIAuthSignUpInput struct {
	EMail    auth.EMail
	Password auth.Password
}

type APIAuthSignUpOutput struct{}

type APIAuthSignInInput struct {
	EMail     auth.EMail
	Password  auth.Password
	PublicKey rsa.PublicKey
	ExpiresIn auth.ExpiresIn
}

type APIAuthSignInOutput struct {
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
}

type APIAuthSignOutInput struct {
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
}

type APIAuthSignOutOutput struct{}

type APIAuthVerifyInput struct {
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
}

type APIAuthVerifyOutput struct{}

type APIAuthGenerateCodeInput struct {
	SessionToken auth.SessionToken
}

type APIAuthGenerateCodeOutput struct {
	Code model.Code
}

type APIAuthRefreshInput struct {
	CodeID       auth.CodeID
	Signature    auth.Signature
	SessionToken auth.SessionToken
}

type APIAuthRefreshOutput struct {
	AuthToken auth.AuthToken
}
