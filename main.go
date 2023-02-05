package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	service "github.com/dzemildupljak/auth-service/internal/core/services"
	"github.com/dzemildupljak/auth-service/internal/db/pgdb"
	"github.com/dzemildupljak/auth-service/internal/handlers/httphdl"
	"github.com/dzemildupljak/auth-service/internal/repositories"
	"github.com/dzemildupljak/auth-service/internal/repositories/persistence"

	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// func init() {

// }

func main() {
	fmt.Println("Hello from auth service!!!")
	ctx := context.Background()

	utils.Load()

	// postgres conn and repo
	pgdbconn := pgdb.DbConnection()
	defer pgdb.CloseDbConnection(pgdbconn)
	persistencerepo := persistence.NewPgRepo(ctx, pgdbconn)
	pgdb.ExecMigrations(pgdbconn)

	// mongo conn and repo
	// dbname := os.Getenv("POSTGRES_DB_AUTH")
	// mngdbconn := mngdb.DbConnection(ctx)
	// mngDB := mngdbconn.Database(dbname)
	// mngdb.ExecMigrations(ctx, mngDB)
	// defer mngdb.DbDisonnection(ctx, mngdbconn)
	// persistencerepo := persistence.NewMngRepo(ctx, mngDB)

	// jwt repo
	jwtrepo := repositories.NewJwtRepo()

	authsrv := service.NewAuthService(ctx, persistencerepo, jwtrepo)
	// authsrv := authservice.NewAuthService(pgrepo, jwtrepo)

	authhdl := httphdl.NewAuthHttpHandler(authsrv)

	usersrv := service.NewUserService(ctx, persistencerepo)
	usrhdl := httphdl.NewUserHttpHandler(usersrv)

	r := mux.NewRouter()

	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("auth healthcheck")
	})

	r.Use(utils.ReqLoggerMiddleware())

	httphdl.AuthRoute(r, *authhdl)
	httphdl.UserRoute(r, *usrhdl)

	appport := os.Getenv("APP_PORT")

	// headers := handlers.AllowedHeaders([]string{"*"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// methods := handlers.AllowedMethods([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000/"})
	ttl := handlers.MaxAge(3600)
	// handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"})

	fmt.Println("ListenAndServe on port :" + appport)
	http.ListenAndServe(":"+appport, handlers.CORS(headers, methods, origins, ttl)(r))
}
