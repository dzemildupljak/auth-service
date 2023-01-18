package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"github.com/dzemildupljak/auth-service/internal/core/ports"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"gopkg.in/validator.v2"
)

type AuthHttpHandler struct {
	service ports.AuthService
}

func NewAuthHttpHandler(srv ports.AuthService) *AuthHttpHandler {

	return &AuthHttpHandler{
		service: srv,
	}
}

func (handler *AuthHttpHandler) Signin(w http.ResponseWriter, r *http.Request) {
	utils.DebugLogger.Println("AuthHttpHandler-Signin start...")

	usr := &domain.UserLogin{}

	err := json.NewDecoder(r.Body).Decode(usr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("decoding error")
		return
	}

	err = validator.Validate(usr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println(err)
		json.NewEncoder(w).Encode("validating error")
		return
	}

	tkns, err := handler.service.Signin(*usr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("authentication error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tkns)
}

func (handler *AuthHttpHandler) Signup(w http.ResponseWriter, r *http.Request) {
	utils.DebugLogger.Println("AuthHttpHandler-Signup start...")

	usr := &domain.SignupUserParams{}

	err := json.NewDecoder(r.Body).Decode(usr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("decoding error")
		return
	}

	err = validator.Validate(usr)
	if err != nil || usr.Password != usr.RPassword {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("validating error")
		return
	}

	err = handler.service.Signup(*usr)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("signup error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (handler *AuthHttpHandler) AuthorizeAccess(w http.ResponseWriter, r *http.Request) {
	token, err := extractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("authentication error")
		return
	}

	err = handler.service.AuthorizeAccess(token)
	if err != nil {
		tokenErrorResponse(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("jwt authentication succeed")
}

func (handler *AuthHttpHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := extractToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("authentication error")
		return
	}
	tkns, err := handler.service.ResetAccesToken(token)
	if err != nil {
		tokenErrorResponse(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tkns)
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	authHeaderContent := strings.Split(authHeader, " ")
	if len(authHeaderContent) != 2 {
		return "", errors.New("token not provided or malformed")
	}

	return authHeaderContent[1], nil
}

func tokenErrorResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode("authentication failed. token not provided or malformed")
	return

}
