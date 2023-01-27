package pgdb

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dzemildupljak/auth-service/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func DbConnection() *gorm.DB {
	var dsn string
	serverconn := os.Getenv("INTERNAL_DATABASE_URL_RENDER")

	if serverconn != "" {
		dsn = serverconn
	} else {
		pguserauth := os.Getenv("POSTGRES_USER")
		pgpassauth := os.Getenv("POSTGRES_PASSWORD")
		pgdbauth := os.Getenv("POSTGRES_DB")
		pgdbhost := os.Getenv("POSTGRES_HOST")

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", pgdbhost, pguserauth, pgpassauth, pgdbauth)

	}

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)

	dbconn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		fmt.Println("Connection failed")
		log.Panic(err)
		return nil
	}

	fmt.Println("Postgres successfully connected...")

	return dbconn
}
func CloseDbConnection(db *gorm.DB) {
	sqlDB, _ := db.DB()

	// Close
	sqlDB.Close()
}
func ExecMigrations(db *gorm.DB) {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	err := db.AutoMigrate(&domain.User{})

	if err != nil {
		fmt.Println("Error migrating postgres")
	} else {
		fmt.Println("Successful migrating postgres")
	}
}
