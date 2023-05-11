package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/core/ports"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/dzemildupljak/auth-service/types"
	"golang.org/x/oauth2"

	"github.com/google/uuid"
)

func authErrorResponse() (types.SigninTokens, error) {
	return types.SigninTokens{
		Access_token:  "",
		Refresh_token: "",
	}, errors.New("authentications failed")
}

type AuthService struct {
	ctx       context.Context
	ssogoogle *oauth2.Config
	prsrepo   ports.PersistenceRepository
	jwtrepo   ports.JwtRepository
	redisrepo ports.RedisRepository
}

func NewAuthService(
	ctx context.Context,
	persrepo ports.PersistenceRepository,
	jwtrepo ports.JwtRepository,
	rdsrepo ports.RedisRepository,
) *AuthService {

	return &AuthService{
		ctx:       ctx,
		prsrepo:   persrepo,
		jwtrepo:   jwtrepo,
		redisrepo: rdsrepo,
	}
}

func (auth *AuthService) Signin(user domain.UserLogin) (types.SigninTokens, error) {
	usr, err := auth.prsrepo.GetUserByMail(user.Email)
	if err != nil {
		fmt.Println("Authservice GetUserByMail failed")
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	correctpwd := utils.ComparePasswords(usr.Password, user.Password)
	if !correctpwd {
		fmt.Println("Authservice ComparePasswords failed")
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	tknresp := types.SigninTokens{}
	tknresp.Access_token, err = auth.jwtrepo.GenerateAccessToken(usr.Id, usr.Role)
	if err != nil {
		fmt.Println("Authservice GenerateAccessToken failed")
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	tknresp.Refresh_token, err = auth.jwtrepo.GenerateRefreshToken(usr.Id, usr.Role)
	if err != nil {
		fmt.Println("Authservice GenerateRefreshToken failed")
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	auth.redisrepo.SetMiddlewareUser(domain.UserMiddleware{
		Id:         usr.Id.String(),
		Email:      usr.Email,
		Username:   usr.Username,
		Isverified: usr.Isverified,
	})

	return tknresp, nil
}

func (auth *AuthService) Signup(user domain.SignupUserParams) error {
	tkhs := utils.GenerateRandomString(64)

	exusr, err := auth.prsrepo.GetUserByMail(user.Email)

	if err == nil && exusr.GoogleId != "" {

		exusr.Password = utils.HashAndSalt(user.Password)
		exusr.Username = user.Username
		exusr.Address = user.Address

		_, err = auth.prsrepo.UpdateRegisterUser(exusr)
		if err != nil {
			fmt.Println("Authservice CreateUpdateRegisterUser failed")
			utils.ErrorLogger.Println(err)
			return err
		}

		return nil
	}

	usr := domain.User{
		Id:         uuid.New(),
		Email:      user.Email,
		Password:   utils.HashAndSalt(user.Password),
		Username:   user.Username,
		Address:    user.Address,
		Name:       user.Name,
		Isverified: false,
		Tokenhash:  []byte(tkhs),
		Role:       "user",
	}

	_, err = auth.prsrepo.CreateRegisterUser(usr)

	if err != nil {
		fmt.Println("Authservice CreateRegisterUser failed")
		utils.ErrorLogger.Println(err)
		return err
	}

	return nil
}

func (auth *AuthService) AuthorizeAccess(acctoken string) error {
	usrid, err := auth.jwtrepo.ValidateAccessToken(acctoken)

	fmt.Println("Authservice ValidateAccessToken for user:", usrid)

	return err
}

func (auth *AuthService) RefreshTokens(reftoken string) (types.SigninTokens, error) {
	usrid, err := auth.jwtrepo.ValidateRefreshToken(reftoken)

	if err != nil {
		fmt.Println("Authservice RefreshTokens failed for user:", usrid)
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	newacctoken, err := auth.jwtrepo.GenerateAccessToken(usrid, "access")

	if err != nil {
		fmt.Println("Authservice GenerateAccessToken failed for user:", usrid)
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	newreftoken, err := auth.jwtrepo.GenerateRefreshToken(usrid, "refresh")

	if err != nil {
		fmt.Println("Authservice GenerateRefreshToken failed for user:", usrid)
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	tknresp := types.SigninTokens{
		Access_token:  newacctoken,
		Refresh_token: newreftoken,
	}

	return tknresp, nil
}
