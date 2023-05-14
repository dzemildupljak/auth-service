package ports

import (
	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/repositories"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(user domain.SignupUserParams) error
	Signin(user domain.UserLogin) (repositories.JwtTokens, error)
	OAuthSignin() (string, error)
	OAuthGoogleCallback(code, state string) (repositories.JwtTokens, error)
	AuthorizeAccess(acctoken string) error
	RefreshTokens(reftoken string) (repositories.JwtTokens, error)
}

type UserService interface {
	GetUserById(usrId uuid.UUID) (domain.User, error)
	GetAllUsers() ([]domain.User, error)
	DeleteUserById(usrId uuid.UUID) error
}
