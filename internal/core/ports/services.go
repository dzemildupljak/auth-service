package ports

import (
	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/types"
	"github.com/google/uuid"
)

type AuthService interface {
	Signin(user domain.UserLogin) (types.SigninTokens, error)
	AuthorizeAccess(acctoken string) error
	ResetTokens(reftoken string) (types.SigninTokens, error)
	Signup(user domain.SignupUserParams) error
}

type UserService interface {
	DeleteUserById(usrId uuid.UUID) error
}
