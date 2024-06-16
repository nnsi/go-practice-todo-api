package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"go-practice-todo/models"
	"go-practice-todo/services"

	"github.com/dgrijalva/jwt-go"
)

type AuthHandler struct {
	service *services.UserService
}

func NewAuthHandler(service *services.UserService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := h.service.HashPassword(user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	if _, err := h.service.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LoginID string `json:"login_id"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.FindByID(req.LoginID)
	if err != nil {
		http.Error(w, "ユーザーIDかパスワードが違います", http.StatusUnauthorized)
		return
	}

	if !h.service.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "ユーザーIDかパスワードが違います", http.StatusUnauthorized)
		return
	}

	token, err := h.service.GenerateToken(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
	WriteJSONResponse(w, token)
}

func (h *AuthHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "無効なアクセスです", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := h.service.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "無効なアクセスです", http.StatusUnauthorized)
			return
		}

		// should not use built-in type string as key for value の解消
		// 独自の型を使ってキーを定義する
		type contextKey string
		const userContextKey contextKey = "user"

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), userContextKey, claims["username"])
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "無効なアクセスです", http.StatusUnauthorized)
		}
	})
}