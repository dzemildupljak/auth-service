package ports

import (
	"context"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
)

type AuthPersistenceRepository interface {
	GetListusers(ctx context.Context) ([]domain.User, error)
	GetUserById(ctx context.Context, id int64) (domain.User, error)
	GetUserByMail(ctx context.Context, mail string) (domain.User, error)
	CreateRegisterUser(ctx context.Context, usr domain.User) error
}

type JwtRepository interface {
	GenerateAccessToken(usrId int64) (string, error)
	ValidateAccessToken(acctoken string) error
	GenerateRefreshToken(usrId int64) (string, error)
	ValidateRefreshToken(reftoken string) (int64, error)
}
