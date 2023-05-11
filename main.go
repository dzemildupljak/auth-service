package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	service "github.com/dzemildupljak/auth-service/internal/core/services"
	"github.com/dzemildupljak/auth-service/internal/db/pgdb"
	"github.com/dzemildupljak/auth-service/internal/handlers/httphdl"
	"github.com/dzemildupljak/auth-service/internal/repositories"
	"github.com/dzemildupljak/auth-service/internal/repositories/persistence"
	"github.com/dzemildupljak/auth-service/internal/utils"
)

func main() {
	fmt.Println("Hello from auth service!!!")
	ctx := context.Background()

	utils.LoadEnv()

	// postgres conn and repo
	pgdbconn := pgdb.DbConnection()
	defer pgdb.CloseDbConnection(pgdbconn)
	pgdb.ExecMigrations(pgdbconn)

	persistencerepo := persistence.NewPgRepo(ctx, pgdbconn)

	// mongo conn and repo
	// dbname := os.Getenv("MONGO_INITDB")
	// mngdbconn := mngdb.DbConnection(ctx)
	// mngDB := mngdbconn.Database(dbname)
	// mngdb.ExecMigrations(ctx, mngDB)
	// defer mngdb.DbDisonnection(ctx, mngdbconn)
	// persistencerepo := persistence.NewMngRepo(ctx, mngDB)

	redisPwd := os.Getenv("REDIS_PWD")
	redislient := redis.NewClient(&redis.Options{
		Addr:     "serviceauthredis:6379",
		Password: redisPwd,
		DB:       0,
	})
	pong, err := redislient.Ping(ctx).Result()
	fmt.Println(pong, err)

	// jwt repo
	jwtrepo := repositories.NewJwtRepo()

	redisrepo := persistence.NewRedisRepo(ctx, redislient)

	authsrv := service.NewAuthService(ctx, persistencerepo, jwtrepo, redisrepo)

	authhdl := httphdl.NewAuthHttpHandler(authsrv)

	usersrv := service.NewUserService(ctx, persistencerepo, redisrepo)
	usrhdl := httphdl.NewUserHttpHandler(usersrv)

	r := mux.NewRouter()

	r.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("auth healthcheck")
	})

	r.Use(utils.ReqLoggerMiddleware())

	httphdl.AuthRoute(r, *authhdl)
	httphdl.UserRoute(r, *usrhdl, persistencerepo, *redisrepo)

	appport := os.Getenv("APP_PORT")
	allowedorigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins(allowedorigins)
	// origins := handlers.AllowedOrigins([]string{"http://localhost:3000/"})
	ttl := handlers.MaxAge(3600)

	fmt.Println("ListenAndServe on port :" + appport)
	http.ListenAndServe(":"+appport, handlers.CORS(headers, methods, origins, ttl)(r))
}
