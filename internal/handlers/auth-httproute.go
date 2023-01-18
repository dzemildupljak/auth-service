package handlers

import (
	"github.com/gorilla/mux"
)

func AuthRoute(r *mux.Router, authhdl AuthHttpHandler) {
	ar := r.PathPrefix("/auth").Subrouter()

	// ar.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode("auth service default encoder")
	// })

	ar.HandleFunc("/", authhdl.AuthorizeAccess)
	ar.HandleFunc("/login", authhdl.Signin).Methods("POST")
	ar.HandleFunc("/signup", authhdl.Signup).Methods("POST")
	ar.HandleFunc("/refresh-token", authhdl.RefreshToken).Methods("GET")
}
