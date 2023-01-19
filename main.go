package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	authservice "github.com/dzemildupljak/auth-service/internal/core/services/auth-service"
	"github.com/dzemildupljak/auth-service/internal/handlers"
	"github.com/dzemildupljak/auth-service/internal/mngdb"
	"github.com/dzemildupljak/auth-service/internal/repositories"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/gorilla/mux"
)

// func init() {

// }

func main() {
	fmt.Println("Hello from auth service!!!")

	utils.Load()

	// postgres conn and repo
	// pgdbconn := pgdb.DbConnection()
	// defer pgdb.CloseDbConnection(pgdbconn)
	// pgrepo := repositories.NewPgRepo(pgdbconn)
	// pgdb.ExecMigrations(dbConn)

	// mongo conn and repo
	ctx := context.Background()
	mngdbconn := mngdb.DbConnection(ctx)
	mngrepo := repositories.NewMngRepo(mngdbconn)

	// jwt repo
	jwtrepo := repositories.NewJwtAuthRepo()

	authsrv := authservice.NewAuthService(mngrepo, jwtrepo)
	// authsrv := authservice.NewAuthService(pgrepo, jwtrepo)
	authhdl := handlers.NewAuthHttpHandler(authsrv)

	r := mux.NewRouter()

	r.Use(utils.ReqLoggerMiddleware())

	handlers.AuthRoute(r, *authhdl)

	appport := os.Getenv("APP_PORT")

	fmt.Println("ListenAndServe on port :" + appport)
	http.ListenAndServe(":"+appport, r)
}
