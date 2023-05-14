package ports

import (
	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/types"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(user domain.SignupUserParams) error
	Signin(user domain.UserLogin) (types.JwtTokens, error)
	OAuthSignin() (string, error)
	OAuthGoogleCallback(code, state string) (types.JwtTokens, error)
	AuthorizeAccess(acctoken string) error
	RefreshTokens(reftoken string) (types.JwtTokens, error)
}

type UserService interface {
	GetUserById(usrId uuid.UUID) (domain.User, error)
	GetAllUsers() ([]domain.User, error)
	DeleteUserById(usrId uuid.UUID) error
}
