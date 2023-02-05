package ports

import (
	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/types"
	"github.com/google/uuid"
)

type AuthService interface {
	Signup(user domain.SignupUserParams) error
	Signin(user domain.UserLogin) (types.SigninTokens, error)
	AuthorizeAccess(acctoken string) error
	RefreshTokens(reftoken string) (types.SigninTokens, error)
}

type UserService interface {
	GetAllUsers() ([]domain.User, error)
	DeleteUserById(usrId uuid.UUID) error
}
