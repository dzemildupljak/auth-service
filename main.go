package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Hello from blog service!!!")

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		"root",
		"postgres",
		"blog-db",
		"5432",
		"blog_service_db")

	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal("connStr", err)
	}
	defer db.Close()

	selectQuer := `SELECT "name" FROM "blogs"`

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

		fmt.Println("BLOG-NAME!!!!!", name)
	}

}
