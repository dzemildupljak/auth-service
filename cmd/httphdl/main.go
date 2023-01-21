package cmd

// import (
// 	"fmt"
// 	"net/http"
// 	"os"

// 	authservice "github.com/dzemildupljak/auth-service/internal/core/services/auth-service"
// 	"github.com/dzemildupljak/auth-service/internal/handlers"
// 	"github.com/dzemildupljak/auth-service/internal/pgdb"
// 	"github.com/dzemildupljak/auth-service/internal/repositories"
// 	"github.com/dzemildupljak/auth-service/internal/utils"
// 	"github.com/gorilla/mux"
// )

// func main() {
// 	fmt.Println("Hello from auth service!!!")

// 	utils.Load()

// 	dbConn := pgdb.DbConnection()
// 	defer pgdb.CloseDbConnection(dbConn)

// 	// pgdb.ExecMigrations(dbConn)

// 	authpgrepo := repositories.NewPgAuthRepo(dbConn)
// 	jwtrepo := repositories.NewJwtRepo()
// 	authsrv := authservice.NewAuthService(authpgrepo, jwtrepo)
// 	authhdl := handlers.NewAuthHttpHandler(authsrv)

// 	r := mux.NewRouter()

// 	r.Use(utils.ReqLoggerMiddleware())

// 	handlers.AuthRoute(r, *authhdl)

// 	appport := os.Getenv("APP_PORT")

// 	fmt.Println("ListenAndServe on port :" + appport)
// 	http.ListenAndServe(":"+appport, r)
// }
