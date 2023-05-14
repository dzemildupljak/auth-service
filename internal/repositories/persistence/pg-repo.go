package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/utils"
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

func (pgrepo *PgRepo) GetUsers() ([]domain.User, error) {
	usrs := []domain.User{}

	err := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Find(&usrs).Error

	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	return usrs, err
}

func (pgrepo *PgRepo) GetUserById(id uuid.UUID) (domain.User, error) {
	usr := domain.User{}

	err := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Where("Id = ?", id).First(&usr).Error
	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	return usr, err
}

func (pgrepo *PgRepo) GetMiddUserById(id uuid.UUID) (domain.UserMiddleware, error) {
	usr := domain.UserMiddleware{}

	err := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Where("Id = ?", id).First(&usr).Error
	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	return usr, err
}

func (pgrepo *PgRepo) GetUserByMail(mail string) (domain.User, error) {
	usr := domain.User{}
	usrQuery := domain.User{Email: mail}

	err := pgrepo.db.WithContext(pgrepo.ctx).First(&usr, usrQuery).Error
	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	return usr, err
}

func (pgrepo *PgRepo) CreateUser(usr domain.User) (domain.User, error) {
	err := validator.Validate(usr)
	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()
	result := pgrepo.db.WithContext(pgrepo.ctx).Create(&usr)

	if result.Error != nil {
		fmt.Println("pg-repo CreateUser failed")
		utils.ErrorLogger.Println(result.Error)
	}

	return usr, result.Error
}

func (pgrepo *PgRepo) CreateOauthUser(usr domain.OauthUserParams) error {
	err := validator.Validate(usr)
	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	usr.CreatedAt = time.Now()
	usr.UpdatedAt = time.Now()

	var cusr domain.User

	err = utils.MapFields(usr, &cusr)
	if err != nil {
		utils.ErrorLogger.Println(err)
		return err
	}

	result := pgrepo.db.WithContext(pgrepo.ctx).Create(&cusr)

	if result.Error != nil {
		utils.ErrorLogger.Println(result.Error)
		return result.Error
	}

	return nil
}

func (pgrepo *PgRepo) UpdateOauthUser(usr domain.OauthUserParams) error {

	err := validator.Validate(usr)
	if err != nil {
		utils.ErrorLogger.Println(err)
		return err
	}

	usr.UpdatedAt = time.Now()

	var cusr domain.User
	err = utils.MapFields(usr, &cusr)
	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	result := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Where("Id = ?", usr.Id).Updates(usr)
	if result.Error != nil {
		utils.ErrorLogger.Println(result.Error)
		return result.Error
	}

	return nil
}

func (pgrepo *PgRepo) UpdateUser(usr domain.User) (domain.User, error) {
	err := validator.Validate(usr)
	if err != nil {
		utils.ErrorLogger.Println(err)
	}

	usr.UpdatedAt = time.Now()
	result := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Where("Id = ?", usr.Id).Updates(
		usr,
	)

	if result.Error != nil {
		utils.ErrorLogger.Println(result.Error)
	}

	return usr, err
}

func (pgrepo *PgRepo) DeleteUserById(id uuid.UUID) error {

	u := domain.User{}

	result := pgrepo.db.WithContext(pgrepo.ctx).Table("users").Where("Id = ?", id).Delete(&u)

	if result.Error != nil || result.RowsAffected == 0 {
		utils.ErrorLogger.Println(result.Error)
		return errors.New("there is no user to delete with given id")
	}

	return result.Error
}
