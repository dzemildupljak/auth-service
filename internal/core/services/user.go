package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/core/ports"
	"github.com/dzemildupljak/auth-service/utils"
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

func (service *UserService) GetAllUsers() ([]domain.User, error) {
	users, err := service.prsrepo.GetUsers()

	if err != nil {
		fmt.Println("Userservice GetUsers falied", err)
		utils.ErrorLogger.Println(err)
	}

	return users, err
}

func (service *UserService) GetUserById(usrId uuid.UUID) (domain.User, error) {
	usr, err := service.prsrepo.GetUserById(usrId)

	if err != nil {
		fmt.Println("Userservice GetUserById falied", err)
		utils.ErrorLogger.Println(err)
	}

	return usr, err
}

func (service *UserService) DeleteUserById(usrId uuid.UUID) error {

	err := service.redisrepo.ClearItemByKey(usrId.String())
	if err != nil {
		fmt.Println("Userservice DeleteUserById redis falied")
		utils.ErrorLogger.Println(err)
		return err
	}

	err = service.prsrepo.DeleteUserById(usrId)
	fmt.Println(err)

	if err != nil {
		fmt.Println("Userservice DeleteUserById falied")
		utils.ErrorLogger.Println(err)
	}

	return err
}
