package service

import (
	"context"

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
	return user.prsrepo.DeleteUserById(usrId)
}

func (user *UserService) GetAllUsers() ([]domain.User, error) {
	return user.prsrepo.GetListusers()
}
