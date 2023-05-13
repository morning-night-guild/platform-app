package usecase

import (
	"context"
	"crypto/rsa"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
)

//go:generate mockgen -source api_auth.go -destination api_auth_mock.go -package usecase

// APIAuth.
type APIAuth interface {
	SignUp(context.Context, APIAuthSignUpInput) (APIAuthSignUpOutput, error)
	SignIn(context.Context, APIAuthSignInInput) (APIAuthSignInOutput, error)
	SignOut(context.Context, APIAuthSignOutInput) (APIAuthSignOutOutput, error)
	SignOutAll(context.Context, APIAuthSignOutAllInput) (APIAuthSignOutAllOutput, error)
	Verify(context.Context, APIAuthVerifyInput) (APIAuthVerifyOutput, error)
	GenerateCode(context.Context, APIAuthGenerateCodeInput) (APIAuthGenerateCodeOutput, error)
	Refresh(context.Context, APIAuthRefreshInput) (APIAuthRefreshOutput, error)
	ChangePassword(context.Context, APIAuthChangePasswordInput) (APIAuthChangePasswordOutput, error)
}

type APIAuthSignUpInput struct {
	EMail    auth.EMail
	Password auth.Password
}

type APIAuthSignUpOutput struct{}

type APIAuthSignInInput struct {
	Secret    auth.Secret
	EMail     auth.EMail
	Password  auth.Password
	PublicKey rsa.PublicKey
	ExpiresIn auth.ExpiresIn
}

type APIAuthSignInOutput struct {
	Auth         model.Auth
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
}

type APIAuthSignOutInput struct {
	UserID    user.ID
	SessionID auth.SessionID
}

type APIAuthSignOutOutput struct{}

type APIAuthSignOutAllInput struct {
	UserID user.ID
}

type APIAuthSignOutAllOutput struct{}

type APIAuthVerifyInput struct {
	UserID    user.ID
	SessionID auth.SessionID
}

type APIAuthVerifyOutput struct{}

type APIAuthGenerateCodeInput struct {
	SessionID auth.SessionID
}

type APIAuthGenerateCodeOutput struct {
	Code model.Code
}

type APIAuthRefreshInput struct {
	CodeID    auth.CodeID
	Signature auth.Signature
	SessionID auth.SessionID
	ExpiresIn auth.ExpiresIn
}

type APIAuthRefreshOutput struct {
	Auth      model.Auth
	AuthToken auth.AuthToken
}

type APIAuthChangePasswordInput struct {
	UserID      user.ID
	Secret      auth.Secret
	PublicKey   rsa.PublicKey
	ExpiresIn   auth.ExpiresIn
	EMail       auth.EMail
	OldPassword auth.Password
	NewPassword auth.Password
}

type APIAuthChangePasswordOutput struct {
	Auth         model.Auth
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
}
