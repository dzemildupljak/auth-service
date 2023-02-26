package service

import (
	"context"
	"fmt"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/core/ports"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/google/uuid"
)

type UserService struct {
	ctx     context.Context
	prsrepo ports.PersistenceRepository
}

func NewUserService(ctx context.Context, authrepo ports.PersistenceRepository) *UserService {
	return &UserService{
		ctx:     ctx,
		prsrepo: authrepo,
	}
}

func (user *UserService) DeleteUserById(usrId uuid.UUID) error {
	err := user.prsrepo.DeleteUserById(usrId)
	if err != nil {
		fmt.Println("Userservice DeleteUserById falied")
		utils.ErrorLogger.Println(err)
	}
	return err
}

func (user *UserService) GetAllUsers() ([]domain.User, error) {
	usrs, err := user.prsrepo.GetListusers()

	if err != nil {
		fmt.Println("Userservice GetListusers falied")
		utils.ErrorLogger.Println(err)
	}
	return usrs, err
}

func (user *UserService) GetUserById(usrId uuid.UUID) (domain.User, error) {
	usr, err := user.prsrepo.GetUserById(usrId)

	if err != nil {
		fmt.Println("Userservice GetUserById falied")
		utils.ErrorLogger.Println(err)
	}
	return usr, err
}
