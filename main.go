package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Hello from auth service!!!")

	dsn := "host=auth-db user=root password=postgres dbname=auth_service_db port=5432 sslmode=disable"

	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Connection failed")
	} else {
		fmt.Println("Connection succeed")
	}

	r := mux.NewRouter()
	ar := r.PathPrefix("/auth").Subrouter()

	ar.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("auth service default encoder")
	})

	ar.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("auth service login encoder")
	})
	http.ListenAndServe(":8004", r)
}
