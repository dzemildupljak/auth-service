package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dzemildupljak/auth-service/internal/domain"
	authservice "github.com/dzemildupljak/auth-service/internal/services/auth-service"
	"github.com/dzemildupljak/auth-service/internal/utils"
)

type AuthHttpHandler struct {
	service authservice.AuthService
}

func NewAuthHttpHandler(srv authservice.AuthService) *AuthHttpHandler {

	return &AuthHttpHandler{
		service: srv,
	}
}

func (handler *AuthHttpHandler) Signin(w http.ResponseWriter, r *http.Request) {
	utils.DebugLogger.Println("AuthHttpHandler-Signin starting..")
	handler.service.Signin(domain.UserLogin{})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("auth service login encoder from handler!!!")
}
