package utils

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnection() *gorm.DB {
	pguserauth := os.Getenv("POSTGRES_USER_AUTH")
	pgpassauth := os.Getenv("POSTGRES_PASSWORD_AUTH")
	pgdbauth := os.Getenv("POSTGRES_DB_AUTH")
	pgdbhost := os.Getenv("POSTGRES_DB_HOST")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", pgdbhost, pguserauth, pgpassauth, pgdbauth)

	dbconn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Connection failed")
		return nil
	} else {
		fmt.Println("Connection succeed")
	}

	return dbconn
}
