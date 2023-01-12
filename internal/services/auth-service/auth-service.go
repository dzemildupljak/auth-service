package authservice

import (
	"context"

	"github.com/dzemildupljak/auth-service/internal/domain"
	"github.com/dzemildupljak/auth-service/internal/ports"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/dzemildupljak/auth-service/types"
)

type AuthService struct {
	repo ports.AuthRepository
}

func NewAuthService(authrepo ports.AuthRepository) *AuthService {
	return &AuthService{
		repo: authrepo,
	}
}

func (auth *AuthService) Signin(user domain.UserLogin) (types.SigninTokens, error) {
	// get user by email and check if exists from adapter(e.g db)
	_, err := auth.repo.GetUserByMail(context.Background(), user.Email)
	if err != nil {
		return types.SigninTokens{}, err
	}

	// compare passwords
	utils.ComparePasswords(utils.HashAndSalt("sifra123"), "sifra123")
	// authenticate

	// generate jwt tokens (access, refresh)

	return types.SigninTokens{
		Access_token:  "",
		Refresh_token: "",
	}, nil
}
