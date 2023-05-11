package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/dzemildupljak/auth-service/types"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (auth *AuthService) OAuthSignin() (string, error) {
	clientId := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectUrl := os.Getenv("GOOGLE_REDIRECT_URL")
	randomString := os.Getenv("GOOGLE_OAUTH_RANDOM_STRING")

	auth.ssogoogle = &oauth2.Config{
		RedirectURL:  redirectUrl,
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid",
		},
		Endpoint: google.Endpoint,
	}
	url := auth.ssogoogle.AuthCodeURL(randomString)

	return url, nil
}

func (auth *AuthService) OAuthGoogleCallback(code, state string) (types.SigninTokens, error) {
	randomString := os.Getenv("GOOGLE_OAUTH_RANDOM_STRING")

	if state != randomString {
		utils.ErrorLogger.Println("wrong state value")
		return types.SigninTokens{}, nil
	}

	token, err := auth.ssogoogle.Exchange(auth.ctx, code)
	if err != nil {
		utils.ErrorLogger.Println("sso google exchange", err)
		return types.SigninTokens{}, nil
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		utils.ErrorLogger.Println("get user info", err)
		return types.SigninTokens{}, nil
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		utils.ErrorLogger.Println("read response body", err)
		return types.SigninTokens{}, nil
	}

	var resdata map[string]interface{}
	err = json.Unmarshal(data, &resdata)
	if err != nil {
		utils.ErrorLogger.Println("unmarshal data", err)
		return types.SigninTokens{}, nil
	}

	usr, err := auth.prsrepo.GetUserByMail(resdata["email"].(string))
	if err == nil {
		acctoken, err := auth.jwtrepo.GenerateAccessToken(usr.Id, usr.Role)
		if err != nil {
			utils.ErrorLogger.Println("callback auth.OAuthSignup GenerateAccessToken", err)
			return types.SigninTokens{}, nil
		}
		refrshtoken, err := auth.jwtrepo.GenerateRefreshToken(usr.Id, usr.Role)
		if err != nil {
			utils.ErrorLogger.Println("callback auth.OAuthSignup GenerateRefreshToken", err)
			return types.SigninTokens{}, nil
		}
		return types.SigninTokens{
			Access_token:  acctoken,
			Refresh_token: refrshtoken,
		}, nil
	}

	tkhs := utils.GenerateRandomString(64)
	fmt.Println("fmt.Println(resdata)", resdata)
	user := domain.SignupOauthUserParams{
		Id:         uuid.New(),
		Email:      resdata["email"].(string),
		Name:       resdata["name"].(string),
		Isverified: resdata["verified_email"].(bool),
		Tokenhash:  []byte(tkhs),
		Role:       "user",
		GoogleId:   resdata["id"].(string),
		Picture:    resdata["picture"].(string),
	}

	fmt.Println("user.Tokenhash", user.Tokenhash)

	err = auth.prsrepo.RegisterOauthUser(user)
	if err != nil {
		utils.ErrorLogger.Println("callback auth.OAuthSignup", err)
		return types.SigninTokens{}, nil
	}
	acctoken, err := auth.jwtrepo.GenerateAccessToken(user.Id, user.Role)
	if err != nil {
		utils.ErrorLogger.Println("callback auth.OAuthSignup GenerateAccessToken", err)
		return types.SigninTokens{}, nil
	}
	refrshtoken, err := auth.jwtrepo.GenerateRefreshToken(user.Id, user.Role)
	if err != nil {
		utils.ErrorLogger.Println("callback auth.OAuthSignup GenerateRefreshToken", err)
		return types.SigninTokens{}, nil
	}

	return types.SigninTokens{
		Access_token:  acctoken,
		Refresh_token: refrshtoken,
	}, nil
}
