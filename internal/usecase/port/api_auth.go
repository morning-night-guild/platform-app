package port

import (
	"crypto/rsa"

	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/auth"
	"github.com/morning-night-guild/platform-app/internal/usecase"
)

type APIAuthSignUpInput struct {
	usecase.Input
	EMail    auth.EMail
	Password auth.Password
}

type APIAuthSignUpOutput struct {
	usecase.Output
}

type APIAuthSignUp interface {
	usecase.Usecase[APIAuthSignUpInput, APIAuthSignUpOutput]
}

type APIAuthSignInInput struct {
	usecase.Input
	EMail     auth.EMail
	Password  auth.Password
	PublicKey rsa.PublicKey
	ExpiresIn auth.ExpiresIn
}

type APIAuthSignInOutput struct {
	usecase.Output
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
}

type APIAuthSignIn interface {
	usecase.Usecase[APIAuthSignInInput, APIAuthSignInOutput]
}

type APIAuthVerifyInput struct {
	usecase.Input
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
}

type APIAuthVerifyOutput struct {
	usecase.Output
}

type APIAuthVerify interface {
	usecase.Usecase[APIAuthVerifyInput, APIAuthVerifyOutput]
}

type APIAuthGenerateCodeInput struct {
	usecase.Input
	SessionToken auth.SessionToken
}

type APIAuthGenerateCodeOutput struct {
	usecase.Output
	Code model.Code
}

type APIAuthGenerateCode interface {
	usecase.Usecase[APIAuthGenerateCodeInput, APIAuthGenerateCodeOutput]
}

type APIAuthRefreshInput struct {
	usecase.Input
	CodeID       auth.CodeID
	Signature    auth.Signature
	SessionToken auth.SessionToken
}

type APIAuthRefreshOutput struct {
	usecase.Output
	AuthToken auth.AuthToken
}

type APIAuthRefresh interface {
	usecase.Usecase[APIAuthRefreshInput, APIAuthRefreshOutput]
}

type APIAuthSignOutInput struct {
	usecase.Input
	AuthToken    auth.AuthToken
	SessionToken auth.SessionToken
}

type APIAuthSignOutOutput struct {
	usecase.Output
}

type APIAuthSignOut interface {
	usecase.Usecase[APIAuthSignOutInput, APIAuthSignOutOutput]
}
