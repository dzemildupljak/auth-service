package authrepo

import (
	"context"

	"github.com/dzemildupljak/auth-service/internal/domain"
	"gorm.io/gorm"
)

type PgAuthRepo struct {
	db gorm.DB
}

func NewPgAuthRepo(dbConn gorm.DB) *PgAuthRepo {
	return &PgAuthRepo{}
}

func (pgarepo *PgAuthRepo) GetListusers(ctx context.Context) ([]domain.User, error) {
	return make([]domain.User, 1), nil
}
func (pgarepo *PgAuthRepo) GetUserById(ctx context.Context, id int64) (domain.User, error) {
	return domain.User{}, nil
}
func (pgarepo *PgAuthRepo) GetUserByMail(ctx context.Context, mail string) (domain.User, error) {
	return domain.User{}, nil
}
