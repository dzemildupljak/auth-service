package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dzemildupljak/auth-service/internal/handlers"
	"github.com/dzemildupljak/auth-service/internal/repositories/authrepo"
	authservice "github.com/dzemildupljak/auth-service/internal/services/auth-service"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Hello from auth service!!!")

	utils.Load()
	pguserauth := os.Getenv("POSTGRES_USER_AUTH")
	pgpassauth := os.Getenv("POSTGRES_PASSWORD_AUTH")
	pgdbauth := os.Getenv("POSTGRES_DB_AUTH")
	pgdbhost := os.Getenv("POSTGRES_DB_HOST")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", pgdbhost, pguserauth, pgpassauth, pgdbauth)

	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Connection failed")
	} else {
		fmt.Println("Connection succeed")
	}

	authrepo := authrepo.NewPgAuthRepo()
	authsrv := authservice.NewAuthService(authrepo)
	authhdl := handlers.NewAuthHttpHandler(*authsrv)

	r := mux.NewRouter()
	ar := r.PathPrefix("/auth").Subrouter()

	ar.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("auth service default encoder")
	})

	ar.HandleFunc("/login", authhdl.Signin)

	// ar.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode("auth service login encoder")
	// })
	http.ListenAndServe(":8004", r)
}
