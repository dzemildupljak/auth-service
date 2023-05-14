package httphdl

import (
	"net/http"

	"github.com/gorilla/mux"
)

func AuthRoute(r *mux.Router, authhdl AuthHttpHandler) {
	ar := r.PathPrefix("/auth").Subrouter()

	ar.HandleFunc("/", authhdl.AuthorizeAccess)
	ar.HandleFunc("/login", authhdl.Signin).Methods("POST")
	ar.HandleFunc("/signup", authhdl.Signup).Methods("POST")
	ar.HandleFunc("/refresh-tokens", authhdl.RefreshToken).Methods("GET")

	o2r := r.PathPrefix("/oauth").Subrouter()

	o2r.HandleFunc("/google/signin", authhdl.GoogleSignin)
	o2r.HandleFunc("/google/callback", authhdl.GoogleCallback)

}

func UserRoute(r *mux.Router, userhdl UserHttpHandler, cachemid HttpCascheMiddlware) {
	ur := r.PathPrefix("/users").Subrouter()

	ur.HandleFunc("", userhdl.ListUser).Methods("GET")
	ur.HandleFunc("/{user_id}", userhdl.GetUserById).Methods("GET")
	ur.HandleFunc("/{user_id}", userhdl.DeleteUserById).Methods("DELETE")
	ur.Use(func(h http.Handler) http.Handler {
		return cachemid.AccessTknMiddleware(h)
	})
}
