package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello from blog service!!!")

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		"root",
		"postgres",
		"auth-db",
		"5432",
		"auth_service_db")

	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("connStr", err)
	}
	defer db.Close()

	selectQuer := `SELECT "name" FROM "users"`

	rows, rErr := db.Query(selectQuer)

	if rErr != nil {
		log.Fatal("select data", rErr)
	}

	defer rows.Close()
	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {
			log.Fatal("rows NEXT", err)
			break
		}

		fmt.Println("USER-NAME!!!!!", name)
	}

	r := mux.NewRouter()
	ar := r.PathPrefix("/auth").Subrouter()

	ar.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
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
