package httphdl

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/dzemildupljak/auth-service/internal/handlers"
// 	"github.com/dzemildupljak/auth-service/internal/pgdb"
// 	"github.com/dzemildupljak/auth-service/internal/repositories/authrepo"
// 	authservice "github.com/dzemildupljak/auth-service/internal/services/auth-service"
// 	"github.com/dzemildupljak/auth-service/internal/utils"
// 	"github.com/gorilla/mux"
// )

// func main() {
// 	fmt.Println("Hello from auth service!!!")

// 	utils.Load()

// 	dbConn := pgdb.DbConnection()

// 	authrepo := authrepo.NewPgAuthRepo(*dbConn)
// 	authsrv := authservice.NewAuthService(authrepo)
// 	authhdl := handlers.NewAuthHttpHandler(*authsrv)

// 	r := mux.NewRouter()

// 	handlers.AuthRoute(r, *authhdl)

// 	http.ListenAndServe(":8004", r)
// }
