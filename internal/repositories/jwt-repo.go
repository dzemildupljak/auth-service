package repositories

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/dzemildupljak/auth-service/types"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtRepo struct {
	config utils.JwtConfigurations
}

func NewJwtRepo() *JwtRepo {
	cnf := utils.NewJwtConfig()

	return &JwtRepo{
		config: cnf,
	}
}

func (jwtrepo *JwtRepo) GenerateTokens(usrId uuid.UUID, urole string) (types.JwtTokens, error) {
	acctoken, err := jwtrepo.GenerateAccessToken(usrId, urole)
	if err != nil {
		return types.JwtTokens{}, err
	}
	reftoken, err := jwtrepo.GenerateRefreshToken(usrId, urole)
	if err != nil {
		return types.JwtTokens{}, err
	}

	return types.JwtTokens{
		Access_token:  acctoken,
		Refresh_token: reftoken,
	}, nil
}

func (jwtrepo *JwtRepo) GenerateAccessToken(usrId uuid.UUID, urole string) (string, error) {
	userID := uuid.UUID.String(usrId)

	tokenType := "access"

	claims := utils.AccessTokenCustomClaims{
		UserId:   userID,
		UserRole: urole,
		KeyType:  tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Second * time.Duration(jwtrepo.config.JwtExpiration),
			).Unix(),
			Issuer: "risc_app.auth.service",
		},
	}

	signBytes, err := os.ReadFile(jwtrepo.config.AccessTokenPrivateKeyPath)

	if err != nil {
		utils.ErrorLogger.Println("unable to read access private key", err)
		fmt.Println("unable to read access private key", err)
		return "", errors.New(
			"could not generate access token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		utils.ErrorLogger.Println("unable to parse private key", "error", err)
		fmt.Println("unable to parse private key", "error", err)
		return "", errors.New(
			"could not generate access token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

func (jwtrepo *JwtRepo) ValidateAccessToken(acctoken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		acctoken,
		&utils.AccessTokenCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				utils.ErrorLogger.Println("Unexpected signing method in auth token")
				fmt.Println("unexpected signing method in auth token")
				return nil, errors.New("unexpected signing method in auth token")
			}
			verifyBytes, err := os.ReadFile(jwtrepo.config.AccessTokenPublicKeyPath)
			if err != nil {
				utils.ErrorLogger.Println("Unable to read public key", "error", err)
				fmt.Println("unable to read public key", "error", err)
				return nil, err
			}

			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
			if err != nil {
				utils.ErrorLogger.Println("Unable to parse public key", "error", err)
				fmt.Println("unable to parse public key", "error", err)
				return nil, err
			}

			return verifyKey, nil
		})

	if err != nil {
		utils.ErrorLogger.Println("Unable to parse claims", "error", err)
		fmt.Println("Unable to parse claims", "error", err)
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*utils.AccessTokenCustomClaims)

	if !ok || !token.Valid || claims.UserId == "" || claims.KeyType != "access" {
		utils.ErrorLogger.Println("invalid token: authentication failed")
		return uuid.Nil, errors.New("invalid token: authentication failed")
	}
	usrid, err := uuid.Parse(claims.UserId)
	if err != nil {
		utils.ErrorLogger.Println("invalid token: authentication failed")
		return uuid.Nil, errors.New("invalid token: authentication failed")
	}

	return usrid, nil
}

func (jwtrepo *JwtRepo) GenerateRefreshToken(userId uuid.UUID, urole string) (string, error) {
	userID := uuid.UUID.String(userId)

	tokenType := "refresh"

	claims := utils.RefreshTokenCustomClaims{
		UserId:   userID,
		KeyType:  tokenType,
		UserRole: urole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Second * time.Duration(jwtrepo.config.JwtExpiration),
			).Unix(),
			Issuer: "risc_app.auth.service",
		},
	}

	signBytes, err := os.ReadFile(jwtrepo.config.RefreshTokenPrivateKeyPath)

	if err != nil {
		utils.ErrorLogger.Println("unable to read refresh private key", err)
		fmt.Println("unable to read refresh private key", err)
		return "", errors.New(
			"could not generate refresh token. please try again later")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		utils.ErrorLogger.Println("unable to parse private key", "error", err)
		fmt.Println("unable to parse private key", "error", err)
		return "", errors.New(
			"could not generate refresh token. please try again later")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}

func (jwtrepo *JwtRepo) ValidateRefreshToken(reftoken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		reftoken,
		&utils.RefreshTokenCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				utils.ErrorLogger.Println("Unexpected signing method in auth token")
				fmt.Println("unexpected signing method in auth token")
				return nil, errors.New("unexpected signing method in auth token")
			}
			verifyBytes, err := os.ReadFile(jwtrepo.config.RefreshTokenPublicKeyPath)
			if err != nil {
				utils.ErrorLogger.Println("Unable to read  refresh public key", "error", err)
				fmt.Println("unable to read public key", "error", err)
				return nil, err
			}

			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
			if err != nil {
				utils.ErrorLogger.Println("Unable to parse refresh public key", "error", err)
				fmt.Println("unable to parse public key", "error", err)
				return nil, err
			}

			return verifyKey, nil
		})

	if err != nil {
		utils.ErrorLogger.Println("Unable to parse claims", "error", err)
		fmt.Println("Unable to parse claims  123", "error", err)
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*utils.RefreshTokenCustomClaims)

	if !ok || !token.Valid || claims.UserId == "" || claims.KeyType != "refresh" {
		utils.ErrorLogger.Println("invalid token: authentication failed")
		return uuid.Nil, errors.New("invalid token: authentication failed")
	}
	usrid, err := uuid.Parse(claims.UserId)
	if err != nil {
		utils.ErrorLogger.Println("invalid token: authentication failed")
		return uuid.Nil, errors.New("invalid token: authentication failed")
	}

	return usrid, nil
}
