package authservice

import (
	"context"

	"github.com/dzemildupljak/auth-service/internal/domain"
	"github.com/dzemildupljak/auth-service/internal/ports"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/dzemildupljak/auth-service/types"
)

type AuthService struct {
	pgrepo ports.AuthRepository
}

func NewAuthService(authrepo ports.AuthRepository) *AuthService {
	return &AuthService{
		pgrepo: authrepo,
	}
}

func (auth *AuthService) Signin(user domain.UserLogin) (types.SigninTokens, error) {
	ctx := context.Background()
	// get user by email and check if exists from adapter(e.g db)
	// usr, err := auth.repo.GetUserByMail(context.Background(), user.Email)

	_, err := auth.pgrepo.GetUserByMail(ctx, user.Email)
	if err != nil {
		return types.SigninTokens{}, err
	}

	// compare passwords
	utils.ComparePasswords(utils.HashAndSalt("sifra123"), "sifra123")
	// authenticate

	// generate jwt tokens (access, refresh)

	return types.SigninTokens{
		Access_token:  "access token value",
		Refresh_token: "refresh token value",
	}, nil
}

func (auth *AuthService) Signup(user domain.SignupUserParams) error {
	ctx := context.Background()

	usr := domain.User{
		Email:      user.Email,
		Password:   utils.HashAndSalt(user.Password),
		Username:   user.Username,
		Address:    user.Address,
		Name:       user.Name,
		Isverified: false,
	}

	auth.pgrepo.CreateRegisterUser(ctx, usr)
	return nil
}
