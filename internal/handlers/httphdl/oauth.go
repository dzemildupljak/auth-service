package httphdl

import (
	"encoding/json"
	"net/http"
)

func (handler *AuthHttpHandler) GoogleSignin(w http.ResponseWriter, r *http.Request) {
	url, err := handler.service.OAuthSignin()
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (handler *AuthHttpHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	usr, err := handler.service.OAuthGoogleCallback(code, state)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(usr)
}
