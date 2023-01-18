package authservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/core/ports"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/dzemildupljak/auth-service/types"
)

func authErrorResponse() (types.SigninTokens, error) {
	return types.SigninTokens{
		Access_token:  "",
		Refresh_token: "",
	}, errors.New("authentications failed")
}

type AuthService struct {
	pgrepo  ports.AuthRepository
	jwtrepo ports.JwtRepository
}

func NewAuthService(authrepo ports.AuthRepository, jwtrepo ports.JwtRepository) *AuthService {
	return &AuthService{
		pgrepo:  authrepo,
		jwtrepo: jwtrepo,
	}
}

func (auth *AuthService) Signin(user domain.UserLogin) (types.SigninTokens, error) {
	ctx := context.Background()

	// get user by email and check if exists from adapter(e.g db)
	usr, err := auth.pgrepo.GetUserByMail(ctx, user.Email)
	if err != nil {
		return types.SigninTokens{}, err
	}

	fmt.Println("fmt.Println(usr)\n\n", usr)

	// compare passwords
	correctpwd := utils.ComparePasswords(usr.Password, user.Password)
	if !correctpwd {
		return authErrorResponse()
	}
	// authenticate

	// generate jwt tokens (access, refresh)
	tknresp := types.SigninTokens{}
	if tknresp.Access_token, err = auth.jwtrepo.GenerateAccessToken(usr.Id); err != nil {
		return authErrorResponse()
	}

	if tknresp.Refresh_token, err = auth.jwtrepo.GenerateRefreshToken(usr.Id); err != nil {
		return authErrorResponse()
	}

	return tknresp, nil
}

func (auth *AuthService) Signup(user domain.SignupUserParams) error {
	ctx := context.Background()
	tkhs := utils.GenerateRandomString(256)

	usr := domain.User{
		Email:      user.Email,
		Password:   utils.HashAndSalt(user.Password),
		Username:   user.Username,
		Address:    user.Address,
		Name:       user.Name,
		Isverified: false,
		Tokenhash:  []byte(tkhs),
	}

	err := auth.pgrepo.CreateRegisterUser(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (auth *AuthService) AuthorizeAccess(acctoken string) error {
	return auth.jwtrepo.ValidateAccessToken(acctoken)
}

func (auth *AuthService) ResetAccesToken(reftoken string) (types.SigninTokens, error) {
	usrid, err := auth.jwtrepo.ValidateRefreshToken(reftoken)
	if err != nil {
		return authErrorResponse()
	}

	newtoken, err := auth.jwtrepo.GenerateAccessToken(usrid)
	if err != nil {
		return authErrorResponse()
	}

	tknresp := types.SigninTokens{
		Access_token:  newtoken,
		Refresh_token: "",
	}

	return tknresp, nil
}
