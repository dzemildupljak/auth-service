package authrepo

import (
	"context"

	"github.com/dzemildupljak/auth-service/internal/domain"
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
	usr := domain.TbUser{}
	pgarepo.db.WithContext(ctx).Table("users").Where("Id = ?", id).First(&usr)

	return domain.User{
		Id:         int64(usr.ID),
		Username:   usr.Username,
		Password:   usr.Password,
		Name:       usr.Name,
		Email:      usr.Email,
		Address:    usr.Address,
		Isverified: usr.Isverified,
	}, nil
}

func (pgarepo *PgAuthRepo) GetUserByMail(ctx context.Context, mail string) (domain.User, error) {
	usr := domain.TbUser{}
	usrQuery := domain.TbUser{Email: mail}
	res := pgarepo.db.First(&usr, usrQuery)

	return domain.User{}, res.Error
}

func (pgarepo *PgAuthRepo) CreateRegisterUser(ctx context.Context, user domain.User) error {
	usr := domain.TbUser{
		Name:       user.Name,
		Username:   user.Username,
		Email:      user.Email,
		Password:   user.Password,
		Address:    user.Address,
		Isverified: user.Isverified,
	}

	res := pgarepo.db.Create(&usr)

	return res.Error
}
