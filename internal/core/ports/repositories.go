package ports

import (
	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/types"
	"github.com/google/uuid"
)

type PersistenceRepository interface {
	GetUsers() ([]domain.User, error)
	GetUserById(id uuid.UUID) (domain.User, error)
	GetUserByMail(mail string) (domain.User, error)
	GetMiddUserById(id uuid.UUID) (domain.UserMiddleware, error)
	DeleteUserById(id uuid.UUID) error
	CreateUser(usr domain.User) (domain.User, error)
	UpdateUser(usr domain.User) (domain.User, error)
	CreateOauthUser(usr domain.OauthUserParams) error
	UpdateOauthUser(usr domain.OauthUserParams) error
}

type JwtRepository interface {
	GenerateTokens(usrId uuid.UUID, urole string) (types.JwtTokens, error)
	GenerateAccessToken(usrId uuid.UUID, urole string) (string, error)
	ValidateAccessToken(acctoken string) (uuid.UUID, error)
	GenerateRefreshToken(userId uuid.UUID, urole string) (string, error)
	ValidateRefreshToken(reftoken string) (uuid.UUID, error)
}

type RedisRepository interface {
	SetMiddlewareUser(usr domain.UserMiddleware) error
	GetMiddlewareUser(usrId string) domain.UserMiddleware
	ClearItemByKey(itemKey string) error
}
