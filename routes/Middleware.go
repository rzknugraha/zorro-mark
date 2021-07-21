package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/context"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

var signMethod = jwt.SigningMethodHS256
var signKey = []byte(viper.GetString("jwt.signature_key"))

// JWTAuthMiddleware middleware func
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Type", "application/json")
		// w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Signing method invalid")
			} else if method != signMethod {
				return nil, fmt.Errorf("Signing method invalid")
			}

			return signKey, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		context.Set(r, "UserInfo", claims)

		next.ServeHTTP(w, r)
	})
}
