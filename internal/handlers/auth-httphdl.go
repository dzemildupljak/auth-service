package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dzemildupljak/auth-service/internal/domain"
	"github.com/dzemildupljak/auth-service/internal/ports"
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

	handler.service.Signup(*usr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
