package ports

import (
	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/google/uuid"
)

type PersistenceRepository interface {
	GetListusers() ([]domain.User, error)
	GetUserById(id uuid.UUID) (domain.User, error)
	DeleteUserById(id uuid.UUID) error
	GetUserByMail(mail string) (domain.User, error)
	CreateRegisterUser(usr domain.User) error
}

type JwtRepository interface {
	GenerateAccessToken(usrId uuid.UUID, urole string) (string, error)
	ValidateAccessToken(acctoken string) (uuid.UUID, error)
	GenerateRefreshToken(userId uuid.UUID, urole string) (string, error)
	ValidateRefreshToken(reftoken string) (uuid.UUID, error)
}
