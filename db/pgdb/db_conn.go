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
	var loglvl int

	workenv, wexsist := os.LookupEnv("WORK_ENVIRONMENT")

	if wexsist && workenv != "local_dev" {
		loglvl = 1
	} else {
		loglvl = 4
	}

	pguserauth := os.Getenv("POSTGRES_USER")
	pgpassauth := os.Getenv("POSTGRES_PASSWORD")
	pgdbauth := os.Getenv("POSTGRES_DB")
	pgdbhost := os.Getenv("POSTGRES_HOST")
	pgdbport := os.Getenv("POSTGRES_PORT")

	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", pgdbhost, pguserauth, pgpassauth, pgdbauth, pgdbport)

	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,             // Slow SQL threshold
			LogLevel:                  logger.LogLevel(loglvl), // Log level
			IgnoreRecordNotFoundError: true,                    // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                    // Disable color
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
