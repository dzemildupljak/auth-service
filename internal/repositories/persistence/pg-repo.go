package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/google/uuid"
	"gopkg.in/validator.v2"
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
	usrs := []domain.User{}

	err := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Find(&usrs).Error

	if err != nil {
		fmt.Println(err)
		utils.ErrorLogger.Println(err)
	}

	return usrs, nil
}

func (pgrepo *PgRepo) GetUserById(id uuid.UUID) (domain.User, error) {
	usr := domain.User{}

	err := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Where("Id = ?", id).First(&usr).Error
	if err != nil {
		fmt.Println(err)
		utils.ErrorLogger.Println(err)
	}

	return usr, nil
}

func (pgrepo *PgRepo) GetUserByMail(mail string) (domain.User, error) {
	usr := domain.User{}
	usrQuery := domain.User{Email: mail}

	err := pgrepo.db.WithContext(pgrepo.ctx).First(&usr, usrQuery).Error
	if err != nil {
		fmt.Println(err)
		utils.ErrorLogger.Println(err)
	}

	return usr, err
}

func (pgrepo *PgRepo) CreateRegisterUser(usr domain.User) error {
	err := validator.Validate(usr)
	if err != nil {
		fmt.Println("persistance:", err)
		utils.ErrorLogger.Println("persistance:", err)
	}

	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()
	err = pgrepo.db.WithContext(pgrepo.ctx).Create(&usr).Error

	if err != nil {
		fmt.Println(err)
		utils.ErrorLogger.Println(err)
	}

	return err
}

func (pgrepo *PgRepo) DeleteUserById(id uuid.UUID) error {

	u := domain.User{}

	err := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Where("Id = ?", id).Delete(&u).Error

	if err != nil {
		fmt.Println(err)
		utils.ErrorLogger.Println(err)
	}

	return err
}
