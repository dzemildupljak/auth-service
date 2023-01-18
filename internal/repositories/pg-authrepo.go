package repositories

import (
	"context"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"gorm.io/gorm"
)

type PgAuthRepo struct {
	db *gorm.DB
}

func NewPgAuthRepo(dbConn *gorm.DB) *PgAuthRepo {
	return &PgAuthRepo{
		db: dbConn,
	}
}

func (pgarepo *PgAuthRepo) GetListusers(ctx context.Context) ([]domain.User, error) {
	return make([]domain.User, 1), nil
}

func (pgarepo *PgAuthRepo) GetUserById(ctx context.Context, id int64) (domain.User, error) {
	usr := domain.User{}
	pgarepo.db.WithContext(ctx).Table("users").Where("Id = ?", id).First(&usr)

	return domain.User{
		Id:         int64(usr.Id),
		Username:   usr.Username,
		Password:   usr.Password,
		Name:       usr.Name,
		Email:      usr.Email,
		Address:    usr.Address,
		Isverified: usr.Isverified,
	}, nil
}

func (pgarepo *PgAuthRepo) GetUserByMail(ctx context.Context, mail string) (domain.User, error) {
	usr := domain.User{}
	usrQuery := domain.User{Email: mail}
	res := pgarepo.db.First(&usr, usrQuery)

	return usr, res.Error
}

func (pgarepo *PgAuthRepo) CreateRegisterUser(ctx context.Context, user domain.User) error {
	res := pgarepo.db.Create(&user)

	return res.Error
}
