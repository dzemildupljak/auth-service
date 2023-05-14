package service

import (
	"encoding/json"
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

func (service *AuthService) OAuthSignin() (string, error) {
	clientId := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectUrl := os.Getenv("GOOGLE_REDIRECT_URL")
	randomString := os.Getenv("GOOGLE_OAUTH_RANDOM_STRING")

	service.ssogoogle = &oauth2.Config{
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
	url := service.ssogoogle.AuthCodeURL(randomString)

	return url, nil
}

func (service *AuthService) OAuthGoogleCallback(code, state string) (types.JwtTokens, error) {
	randomString := os.Getenv("GOOGLE_OAUTH_RANDOM_STRING")

	if state != randomString {
		utils.ErrorLogger.Println("wrong state value")
		return types.JwtTokens{}, nil
	}

	token, err := service.ssogoogle.Exchange(service.ctx, code)
	if err != nil {
		utils.ErrorLogger.Println("sso google exchange", err)
		return authErrorResponse()
	}

	googleres, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		utils.ErrorLogger.Println("get user info", err)
		return authErrorResponse()
	}

	defer googleres.Body.Close()

	resbody, err := io.ReadAll(googleres.Body)
	if err != nil {
		utils.ErrorLogger.Println("read response body", err)
		return authErrorResponse()
	}

	var userdata map[string]interface{}
	err = json.Unmarshal(resbody, &userdata)
	if err != nil {
		utils.ErrorLogger.Println("unmarshal data", err)
		return authErrorResponse()
	}

	usr, err := service.prsrepo.GetUserByMail(userdata["email"].(string))

	if err == nil {
		if usr.GoogleId != "" && usr.GoogleId != userdata["id"].(string) {
			utils.ErrorLogger.Println("wrong google id", err)
			return authErrorResponse()
		}
		ouser := domain.OauthUserParams{
			Id:         usr.Id,
			Name:       userdata["name"].(string),
			Isverified: userdata["verified_email"].(bool),
			GoogleId:   userdata["id"].(string),
			Picture:    userdata["picture"].(string),
			Tokenhash:  usr.Tokenhash,
			Role:       usr.Role,
			Email:      usr.Email,
		}
		err = service.prsrepo.UpdateOauthUser(ouser)

		if err != nil {
			utils.ErrorLogger.Println("callback service.prsrepo.UpdateOauthUser(usr)", err)
			return authErrorResponse()
		}

		jwttokens, err := service.jwtrepo.GenerateTokens(usr.Id, usr.Role)
		if err != nil {
			utils.ErrorLogger.Println("callback service.jwtrepo.GenerateTokens", err)
			return authErrorResponse()
		}
		return jwttokens, nil
	}

	tkhs := utils.GenerateRandomString(64)

	user := domain.OauthUserParams{
		Id:         uuid.New(),
		Email:      userdata["email"].(string),
		Name:       userdata["name"].(string),
		Isverified: userdata["verified_email"].(bool),
		Tokenhash:  []byte(tkhs),
		Role:       "user",
		GoogleId:   userdata["id"].(string),
		Picture:    userdata["picture"].(string),
	}

	err = service.prsrepo.CreateOauthUser(user)

	if err != nil {
		utils.ErrorLogger.Println("callback auth.OAuthSignup", err)
		return authErrorResponse()
	}

	jwttokens, err := service.jwtrepo.GenerateTokens(user.Id, user.Role)
	if err != nil {
		utils.ErrorLogger.Println("callback service.jwtrepo.GenerateTokens", err)
		return authErrorResponse()
	}
	return jwttokens, nil
}
