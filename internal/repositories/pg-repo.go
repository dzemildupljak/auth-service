package repositories

import (
	"context"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"gorm.io/gorm"
)

type PgRepo struct {
	db *gorm.DB
}

func NewPgRepo(dbConn *gorm.DB) *PgRepo {
	return &PgRepo{
		db: dbConn,
	}
}

func (pgrepo *PgRepo) GetListusers(ctx context.Context) ([]domain.User, error) {
	return make([]domain.User, 1), nil
}

func (pgrepo *PgRepo) GetUserById(ctx context.Context, id int64) (domain.User, error) {
	usr := domain.User{}
	pgrepo.db.WithContext(ctx).Table("users").Where("Id = ?", id).First(&usr)

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

func (pgrepo *PgRepo) GetUserByMail(ctx context.Context, mail string) (domain.User, error) {
	usr := domain.User{}
	usrQuery := domain.User{Email: mail}
	res := pgrepo.db.First(&usr, usrQuery)

	return usr, res.Error
}

func (pgrepo *PgRepo) CreateRegisterUser(ctx context.Context, user domain.User) error {
	res := pgrepo.db.Create(&user)

	return res.Error
}
