package service

import (
	"context"
	"errors"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/core/ports"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/dzemildupljak/auth-service/types"
	"github.com/google/uuid"
)

func authErrorResponse() (types.SigninTokens, error) {
	return types.SigninTokens{
		Access_token:  "",
		Refresh_token: "",
	}, errors.New("authentications failed")
}

type AuthService struct {
	ctx     context.Context
	prsrepo ports.PersistenceRepository
	jwtrepo ports.JwtRepository
}

func NewAuthService(ctx context.Context, persrepo ports.PersistenceRepository, jwtrepo ports.JwtRepository) *AuthService {
	return &AuthService{
		ctx:     ctx,
		prsrepo: persrepo,
		jwtrepo: jwtrepo,
	}
}

func (auth *AuthService) Signin(user domain.UserLogin) (types.SigninTokens, error) {

	// get user by email and check if exists from adapter(e.g db)
	usr, err := auth.prsrepo.GetUserByMail(user.Email)

	if err != nil {
		return types.SigninTokens{}, err
	}

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
	tkhs := utils.GenerateRandomString(64)

	usr := domain.User{
		Id:         uuid.New(),
		Email:      user.Email,
		Password:   utils.HashAndSalt(user.Password),
		Username:   user.Username,
		Address:    user.Address,
		Name:       user.Name,
		Isverified: false,
		Tokenhash:  []byte(tkhs),
	}

	err := auth.prsrepo.CreateRegisterUser(usr)
	if err != nil {
		utils.ErrorLogger.Println(err)
		return err
	}
	return nil
}

func (auth *AuthService) AuthorizeAccess(acctoken string) error {
	_, err := auth.jwtrepo.ValidateAccessToken(acctoken)
	return err
}

func (auth *AuthService) ResetTokens(reftoken string) (types.SigninTokens, error) {
	usrid, err := auth.jwtrepo.ValidateRefreshToken(reftoken)
	if err != nil {
		return authErrorResponse()
	}

	newacctoken, err := auth.jwtrepo.GenerateAccessToken(usrid)
	if err != nil {
		return authErrorResponse()
	}
	newreftoken, err := auth.jwtrepo.GenerateRefreshToken(usrid)
	if err != nil {
		return authErrorResponse()
	}

	tknresp := types.SigninTokens{
		Access_token:  newacctoken,
		Refresh_token: newreftoken,
	}

	return tknresp, nil
}
