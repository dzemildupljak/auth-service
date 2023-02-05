package repositories

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JwtRepo struct {
	config utils.JwtConfigurations
}

func NewJwtRepo() *JwtRepo {
	conf := utils.NewJwtConfig()

	return &JwtRepo{
		config: conf,
	}
}

func (jwtrepo *JwtRepo) GenerateAccessToken(usrId uuid.UUID) (string, error) {
	userID := uuid.UUID.String(usrId)

	tokenType := "access"
	userRole := "user"

	claims := utils.AccessTokenCustomClaims{
		UserId:   userID,
		UserRole: userRole,
		KeyType:  tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Second * time.Duration(jwtrepo.config.JwtExpiration),
			).Unix(),
			Issuer:    "risc_app.auth.service",
			IssuedAt:  time.Now().UnixMilli(),
			NotBefore: time.Now().UnixMilli(),
		},
	}

	signBytes, err := os.ReadFile(jwtrepo.config.AccessTokenPrivateKeyPath)

	if err != nil {
		utils.ErrorLogger.Println("unable to read access private key", err)
		fmt.Println("unable to read access private key", err)
		return "", errors.New(
			"could not generate access token. please try again later 1")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		utils.ErrorLogger.Println("unable to parse private key", "error", err)
		fmt.Println("unable to parse private key", "error", err)
		return "", errors.New(
			"could not generate access token. please try again later 2")
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
		return uuid.Nil, errors.New("invalid token: authentication failed")
	}
	usrid, err := uuid.Parse(claims.UserId)
	if err != nil {
		return uuid.Nil, errors.New("invalid token: authentication failed")
	}

	return usrid, nil
}

func (jwtrepo *JwtRepo) GenerateRefreshToken(userId uuid.UUID) (string, error) {
	usrId := uuid.UUID.String(userId)
	cusKey := utils.GenerateCustomKey(usrId, "asdadsads")
	tokenType := "refresh"

	claims := utils.RefreshTokenCustomClaims{
		UserId:    usrId,
		CustomKey: cusKey,
		KeyType:   tokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				24 * time.Duration(jwtrepo.config.JwtRefreshExpiration),
			).Unix(),
			Issuer:    "risc_app.auth.service",
			IssuedAt:  time.Now().UnixMilli(),
			NotBefore: time.Now().UnixMilli(),
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
				utils.ErrorLogger.Println("unexpected signing method in auth token")
				fmt.Println("unexpected signing method in auth token")
				return nil, errors.New("unexpected signing method in auth token")
			}
			verifyBytes, err := os.ReadFile(jwtrepo.config.RefreshTokenPublicKeyPath)
			if err != nil {
				utils.ErrorLogger.Println("unable to read public key", "error", err)
				fmt.Println("unable to read public key", "error", err)
				return nil, err
			}

			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
			if err != nil {
				utils.ErrorLogger.Println("unable to parse public key", "error", err)
				fmt.Println("unable to parse public key", "error", err)
				return nil, err
			}

			return verifyKey, nil
		})

	if err != nil {
		utils.ErrorLogger.Println("unable to parse claims", "error", err)
		fmt.Println("unable to parse claims", "error", err)
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*utils.RefreshTokenCustomClaims)
	usrId, err := uuid.Parse(claims.UserId)

	if !ok || !token.Valid || claims.UserId == "" || claims.KeyType != "refresh" || err != nil {
		utils.ErrorLogger.Println("could not extract claims from token")
		fmt.Println("could not extract claims from token")
		return uuid.Nil, errors.New("invalid token: authentication failed")
	}

	return usrId, nil
}
