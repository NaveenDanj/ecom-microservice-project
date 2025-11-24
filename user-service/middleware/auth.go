package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"user-service/config"
	"user-service/models"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const UserIDKey ctxKey = "userID"
const UserRoleKey ctxKey = "userRole"

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		fmt.Printf("Token String: %s\n", tokenString)
		fmt.Printf("Key : %s\n", os.Getenv("JWT_SECRET"))

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		id := claims["id"].(string)
		role := claims["role"].(string)

		fmt.Println("Authenticated User ID:", id, "Role:", role)

		var user models.User

		if err := config.DB.Model(&models.User{}).Where("id = ?", id).First(&user).Error; err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if role != "user" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if user.IsActive == false {
			http.Error(w, "Account is inactive", http.StatusForbidden)
			return
		}

		if user.IsVerified == false {
			http.Error(w, "Account is not verified", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, id)
		ctx = context.WithValue(ctx, UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(UserIDKey)
	s, ok := v.(string)
	return s, ok
}
