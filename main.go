package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dzemildupljak/auth-service/internal/handlers"
	"github.com/dzemildupljak/auth-service/internal/pgdb"
	"github.com/dzemildupljak/auth-service/internal/repositories/authrepo"
	authservice "github.com/dzemildupljak/auth-service/internal/services/auth-service"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Hello from auth service!!!")

	utils.Load()

	dbConn := pgdb.DbConnection()
	defer pgdb.CloseDbConnection(dbConn)

	// pgdb.ExecMigrations(dbConn)

	authpgrepo := authrepo.NewPgAuthRepo(dbConn)
	authsrv := authservice.NewAuthService(authpgrepo)
	authhdl := handlers.NewAuthHttpHandler(authsrv)

	r := mux.NewRouter()

	handlers.AuthRoute(r, *authhdl)

	appport := os.Getenv("APP_PORT")

	fmt.Println("ListenAndServe on port :" + appport)
	http.ListenAndServe(":"+appport, r)
}
