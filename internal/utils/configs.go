package utils

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt"
)

type AccessTokenCustomClaims struct {
	UserId   string `json:"uid"`
	UserRole string `json:"urole"`
	KeyType  string `json:"type"`
	jwt.StandardClaims
}
type RefreshTokenCustomClaims struct {
	UserId  string `json:"uid"`
	KeyType string `json:"type"`
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

func NewJwtConfig() JwtConfigurations {
	curDir, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	return JwtConfigurations{
		AccessTokenPrivateKeyPath:  curDir + "/access-private.pem",
		AccessTokenPublicKeyPath:   curDir + "/access-public.pem",
		RefreshTokenPrivateKeyPath: curDir + "/refresh-private.pem",
		RefreshTokenPublicKeyPath:  curDir + "/refresh-public.pem",
		JwtExpiration:              60,  // seconds
		JwtRefreshExpiration:       360, // seconds
	}
}
