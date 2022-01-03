package main

import (
	"database/sql"
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

	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello from auth")
	})
	http.ListenAndServe(":8004", r)

}
