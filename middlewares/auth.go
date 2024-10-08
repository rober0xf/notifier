package middlewares

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"goapi/config"
	"net/http"
	"strings"
)

var jwtKey = []byte(config.JwtKey)

// custom jwt message. http.error style
func jwtErrorMessage(w http.ResponseWriter, err string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jwtErrorMessage(w, "missing or malformed jwt", http.StatusBadRequest)
			return
		}

		// remove the bearer from the header to keep only the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			jwtErrorMessage(w, "missing or malformed jwt", http.StatusBadRequest)
			return
		}

		// extract the information from the token
		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return []byte(config.JwtKey), nil
		})
		if err != nil {
			jwtErrorMessage(w, "invalid jwt", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			jwtErrorMessage(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// if both the token and authentication is valid, the middleware calls the next handler
		next.ServeHTTP(w, r)
	})
}

func ProtectedTest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("you can have access"))
}
