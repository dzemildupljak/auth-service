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

func authErrorResponse() (types.JwtTokens, error) {
	return types.JwtTokens{
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

func (service *AuthService) Signup(user domain.SignupUserParams) error {
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
		Role:       "user",
	}

	_, err := service.prsrepo.CreateUser(usr)

	if err != nil {
		fmt.Println("Authservice CreateUser failed")
		utils.ErrorLogger.Println(err)
		return err
	}

	return nil
}

func (service *AuthService) Signin(user domain.UserLogin) (types.JwtTokens, error) {
	usr, err := service.prsrepo.GetUserByMail(user.Email)
	if err != nil {
		fmt.Println("Authservice GetUserByMail failed")
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	okpwd := utils.ComparePasswords(usr.Password, user.Password)
	if !okpwd {
		fmt.Println("Authservice ComparePasswords failed")
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	tokens, err := service.jwtrepo.GenerateTokens(usr.Id, usr.Role)
	if err != nil {
		fmt.Println("Authservice GenerateTokens failed", err)
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	service.redisrepo.SetMiddlewareUser(domain.UserMiddleware{
		Id:         usr.Id.String(),
		Email:      usr.Email,
		Username:   usr.Username,
		Isverified: usr.Isverified,
	})

	return tokens, nil
}

func (service *AuthService) AuthorizeAccess(acctoken string) error {
	_, err := service.jwtrepo.ValidateAccessToken(acctoken)

	return err
}

func (service *AuthService) RefreshTokens(reftoken string) (types.JwtTokens, error) {
	usrid, err := service.jwtrepo.ValidateRefreshToken(reftoken)
	if err != nil {
		fmt.Println("Authservice RefreshTokens failed for user:", usrid)
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	user, err := service.prsrepo.GetUserById(usrid)
	if err != nil {
		fmt.Println("Authservice RefreshTokens GetUserById failed for user:", usrid)
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	tokens, err := service.jwtrepo.GenerateTokens(usrid, user.Role)
	if err != nil {
		fmt.Println("Authservice GenerateTokens failed for user:", usrid)
		utils.ErrorLogger.Println(err)
		return authErrorResponse()
	}

	return tokens, nil
}
