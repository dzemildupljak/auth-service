package repositories

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/golang-jwt/jwt"
)

// AccessTokenCustomClaims specifies the claims for access token
type AccessTokenCustomClaims struct {
	UserID   string
	UserRole string
	KeyType  string
	jwt.StandardClaims
}
type RefreshTokenCustomClaims struct {
	UserID    string
	CustomKey string
	KeyType   string
	jwt.StandardClaims
}

type JwtConfigurations struct {
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	JwtExpiration              int
	JwtRefreshExpiration       int
	MailVerifTemplateID        string
	PassResetTemplateID        string
}

type JwtAuthRepo struct {
	config JwtConfigurations
}

func NewJwtAuthRepo() *JwtAuthRepo {
	curDir, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	return &JwtAuthRepo{
		config: JwtConfigurations{
			AccessTokenPrivateKeyPath:  curDir + "/access-private.pem",
			AccessTokenPublicKeyPath:   curDir + "/access-public.pem",
			RefreshTokenPrivateKeyPath: curDir + "/refresh-private.pem",
			RefreshTokenPublicKeyPath:  curDir + "/refresh-public.pem",
			JwtExpiration:              60,  // seconds
			JwtRefreshExpiration:       360, // seconds
		},
	}
}

func (jwtrepo *JwtAuthRepo) GenerateAccessToken(usrId int64) (string, error) {
	userID := strconv.FormatInt(usrId, 10)
	tokenType := "access"
	userRole := "user"

	claims := AccessTokenCustomClaims{
		userID,
		userRole,
		tokenType,
		jwt.StandardClaims{
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

func (jwtrepo *JwtAuthRepo) ValidateAccessToken(acctoken string) error {
	token, err := jwt.ParseWithClaims(
		acctoken,
		&AccessTokenCustomClaims{},
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
		return err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)

	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "access" {
		return errors.New("invalid token: authentication failed")
	}

	return nil
}

func (jwtrepo *JwtAuthRepo) GenerateRefreshToken(userId int64) (string, error) {
	usrId := strconv.FormatInt(userId, 10)
	cusKey := utils.GenerateCustomKey(usrId, "asdadsads")
	tokenType := "refresh"

	claims := RefreshTokenCustomClaims{
		usrId,
		cusKey,
		tokenType,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				24 * time.Hour * time.Duration(jwtrepo.config.JwtRefreshExpiration),
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

func (jwtrepo *JwtAuthRepo) ValidateRefreshToken(reftoken string) (int64, error) {

	token, err := jwt.ParseWithClaims(
		reftoken,
		&RefreshTokenCustomClaims{},
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
		return 0, err
	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	usrId, err := strconv.ParseInt(claims.UserID, 10, 64)

	if !ok || !token.Valid || claims.UserID == "" || claims.KeyType != "refresh" || err != nil {
		utils.ErrorLogger.Println("could not extract claims from token")
		fmt.Println("could not extract claims from token")
		return 0, errors.New("invalid token: authentication failed")
	}

	return usrId, nil
}
