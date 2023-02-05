package httphdl

import (
	"github.com/gorilla/mux"
)

func AuthRoute(r *mux.Router, authhdl AuthHttpHandler) {
	ar := r.PathPrefix("/auth").Subrouter()

	ar.HandleFunc("/", authhdl.AuthorizeAccess)
	ar.HandleFunc("/login", authhdl.Signin).Methods("POST")
	ar.HandleFunc("/signup", authhdl.Signup).Methods("POST")
	ar.HandleFunc("/refresh-tokens", authhdl.RefreshToken).Methods("GET")
}

func UserRoute(r *mux.Router, userhdl UserHttpHandler) {
	ur := r.PathPrefix("/user").Subrouter()

	ur.HandleFunc("/users", userhdl.ListUser).Methods("GET")
	ur.HandleFunc("/{user_id}", userhdl.DeleteUserById).Methods("DELETE")
}
