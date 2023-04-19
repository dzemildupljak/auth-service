package httphdl

import (
	"encoding/json"
	"net/http"

	"github.com/dzemildupljak/auth-service/internal/core/ports"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHttpHandler struct {
	service ports.UserService
}

func NewUserHttpHandler(srv ports.UserService) *UserHttpHandler {
	return &UserHttpHandler{
		service: srv,
	}
}

func (handler *UserHttpHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	usrs, err := handler.service.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user list retrieve failed")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{
		Payload: usrs,
		Status:  true,
	})
}

func (handler *UserHttpHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usrid, err := uuid.Parse(params["user_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user not found")
		return
	}

	usr, err := handler.service.GetUserById(usrid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ResponsePayload{
		Payload: usr,
		Status:  true,
	})

}
func (handler *UserHttpHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	usrid, err := uuid.Parse(params["user_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user delete failed")
		return
	}

	err = handler.service.DeleteUserById(usrid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("user delete failed")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("deleted")
}
