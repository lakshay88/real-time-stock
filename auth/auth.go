package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lakshay88/real-time-stock/internal/user/models"
)

type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	UserEmail string `json:"email"`
	jwt.StandardClaims
}

var jwtSecret = []byte("your_secret_key")

func GenerateJWTToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID:    user.ID,
		Username:  user.Username,
		UserEmail: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		splitedValue := strings.Split(authHeader, " ")
		if len(splitedValue) != 2 || splitedValue[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(splitedValue[1], claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", claims.UserID)
		ctx = context.WithValue(ctx, "email", claims.UserEmail)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
