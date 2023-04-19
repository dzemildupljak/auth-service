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
	ctx       context.Context
	prsrepo   ports.PersistenceRepository
	redisrepo ports.RedisRepository
}

func NewUserService(
	ctx context.Context,
	prsrepo ports.PersistenceRepository,
	redisrepo ports.RedisRepository,
) *UserService {
	return &UserService{
		ctx:       ctx,
		prsrepo:   prsrepo,
		redisrepo: redisrepo,
	}
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

func (user *UserService) DeleteUserById(usrId uuid.UUID) error {

	err := user.redisrepo.ClearItemByKey(usrId.String())
	if err != nil {
		fmt.Println("Userservice DeleteUserById redis falied")
		utils.ErrorLogger.Println(err)
		return err
	}

	err = user.prsrepo.DeleteUserById(usrId)
	fmt.Println(err)

	if err != nil {
		fmt.Println("Userservice DeleteUserById falied")
		utils.ErrorLogger.Println(err)
	}

	return err
}
