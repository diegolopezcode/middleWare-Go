package middleware

import (
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

var JWT_SECRET string

func setJWTSecret(secret string) {
	JWT_SECRET = secret
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(
				w, "Authorization Header required", http.StatusUnauthorized,
			)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer")
		if tokenString == authHeader {
			http.Error(w, "Invalid Authorization Format", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("user_id", claims["user_id"].(string))
		next.ServeHTTP(w, r)
	})
}
