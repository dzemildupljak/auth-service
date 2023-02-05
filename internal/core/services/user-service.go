package service

import (
	"context"
	"fmt"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/core/ports"
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
	}
	return err
}

func (user *UserService) GetAllUsers() ([]domain.User, error) {
	usrs, err := user.prsrepo.GetListusers()

	if err != nil {
		fmt.Println("Userservice GetListusers falied")
	}
	return usrs, err
}
