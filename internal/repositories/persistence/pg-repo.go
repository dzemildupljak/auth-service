package persistence

import (
	"context"
	"time"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PgRepo struct {
	db  *gorm.DB
	ctx context.Context
}

func NewPgRepo(ctx context.Context, dbConn *gorm.DB) *PgRepo {
	return &PgRepo{
		db:  dbConn,
		ctx: ctx,
	}
}

func (pgrepo *PgRepo) GetListusers() ([]domain.User, error) {
	return make([]domain.User, 1), nil
}

func (pgrepo *PgRepo) GetUserById(id uuid.UUID) (domain.User, error) {
	usr := domain.User{}
	err := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Where("Id = ?", id).First(&usr).Error
	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	return usr, nil
}

func (pgrepo *PgRepo) GetUserByMail(mail string) (domain.User, error) {
	usr := domain.User{}
	usrQuery := domain.User{Email: mail}
	res := pgrepo.db.WithContext(pgrepo.ctx).First(&usr, usrQuery)

	if res.Error != nil {
		utils.ErrorLogger.Println(res.Error)
	}

	return usr, res.Error
}

func (pgrepo *PgRepo) CreateRegisterUser(usr domain.User) error {
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()
	res := pgrepo.db.WithContext(pgrepo.ctx).Create(&usr)

	if res.Error != nil {
		utils.ErrorLogger.Println(res.Error)
	}

	return res.Error
}

func (pgrepo *PgRepo) DeleteUserById(id uuid.UUID) error {
	return nil
}
