package httphdl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/dzemildupljak/auth-service/internal/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func (handler *AuthHttpHandler) GoogleSignin(w http.ResponseWriter, r *http.Request) {
	clientId := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	redirectUrl := os.Getenv("GOOGLE_REDIRECT_URL")
	randomString := os.Getenv("GOOGLE_OAUTH_RANDOM_STRING")

	handler.ssogoogle = &oauth2.Config{
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
	url := handler.ssogoogle.AuthCodeURL(randomString)
	fmt.Println(url)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (handler *AuthHttpHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")
	randomString := os.Getenv("GOOGLE_OAUTH_RANDOM_STRING")

	if state != randomString {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := handler.ssogoogle.Exchange(r.Context(), code)
	if err != nil {
		utils.ErrorLogger.Println("sso google exchange", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		utils.ErrorLogger.Println("get user info", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		utils.ErrorLogger.Println("read response body", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var resdata map[string]interface{}
	err = json.Unmarshal(data, &resdata)
	if err != nil {
		utils.ErrorLogger.Println("unmarshal data", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resdata)
}

// {
// 	"email": "dzemildupljak4795@gmail.com",
// 	"family_name": "dupljak",
// 	"given_name": "dzemil",
// 	"id": "117861098009647898666",
// 	"locale": "sr",
// 	"name": "dzemil dupljak",
// 	"picture": "https://lh3.googleusercontent.com/a/AGNmyxatZvTwmO4yB80PZDr7GVBqdelhAazxF8t_gSMOZw=s96-c",
// 	"verified_email": true
//   }
