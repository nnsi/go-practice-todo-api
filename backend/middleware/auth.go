package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-practice-todo/services"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

const userContextKey contextKey = "user_id"

func AuthMiddleware(userService *services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorizationヘッダーが不足しています", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := userService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "無効なトークンです", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID, ok := claims["user_id"].(string)
				if !ok {
					http.Error(w, "トークンにユーザー名が含まれていません", http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), userContextKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				http.Error(w, "無効なトークンです", http.StatusUnauthorized)
			}
		})
	}
}

func GetUserFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(userContextKey).(string); ok {
		return userID
	}
	return ""
}
