package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func AuthRoute(r *mux.Router, authhdl AuthHttpHandler) {
	ar := r.PathPrefix("/auth").Subrouter()

	ar.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("auth service default encoder")
	})

	ar.HandleFunc("/login", authhdl.Signin)
}
