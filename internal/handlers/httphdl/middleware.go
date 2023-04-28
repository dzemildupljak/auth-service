package httphdl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dzemildupljak/auth-service/internal/core/ports"
	"github.com/dzemildupljak/auth-service/internal/repositories/persistence"
	"github.com/dzemildupljak/auth-service/internal/utils"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func AccTknMiddleware(next http.Handler, persrepo ports.PersistenceRepository, redis persistence.RedisRepo) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		acctoken, err := extractToken(r)
		if err != nil {
			TokenErrorResponse(w)
			return
		}

		config := utils.NewJwtConfig()
		token, err := jwt.ParseWithClaims(
			acctoken,
			&utils.AccessTokenCustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					utils.ErrorLogger.Println("Unexpected signing method in auth token")
					fmt.Println("unexpected signing method in auth token")
					return nil, errors.New("unexpected signing method in auth token")
				}
				verifyBytes, err := os.ReadFile(config.AccessTokenPublicKeyPath)
				if err != nil {
					utils.ErrorLogger.Println("Unable to read public key", "error", err)
					fmt.Println("unable to read public key", "error", err)
					return nil, err
				}

				verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
				if err != nil {
					utils.ErrorLogger.Println("Unable to parse public key", "error", err)
					fmt.Println("unable to parse public key", "error", err)
					return nil, err
				}

				return verifyKey, nil
			})

		if err != nil {
			fmt.Println(err)
			TokenErrorResponse(w)
			return
		}

		claims, ok := token.Claims.(*utils.AccessTokenCustomClaims)
		if !ok || !token.Valid || claims.UserId == "" || claims.KeyType != "access" {
			TokenErrorResponse(w)
			return
		}

		_, err = uuid.Parse(claims.UserId)
		if err != nil {
			TokenErrorResponse(w)
			return
		}

		midusr := redis.GetMiddlewareUser(claims.UserId)

		if midusr.Id != claims.UserId {
			usrMidd, err := persrepo.GetMiddUserById(uuid.MustParse(claims.UserId))
			if err != nil || usrMidd.Id != claims.UserId {
				TokenErrorResponse(w)
				return
			} else {
				redis.SetMiddlewareUser(usrMidd)
			}
		} else if !midusr.Isverified {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("authentication failed, user is not verified")
			return
		}

		fmt.Println("END OF MIDDLEWARE")

		next.ServeHTTP(w, r)
	})
}
