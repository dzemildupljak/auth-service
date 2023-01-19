package ports

import (
	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/types"
)

type AuthService interface {
	Signin(user domain.UserLogin) (types.SigninTokens, error)
	AuthorizeAccess(acctoken string) error
	ResetAccesToken(acctoken string) (types.SigninTokens, error)
	Signup(user domain.SignupUserParams) error
}
