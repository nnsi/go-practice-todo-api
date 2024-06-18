package handlers

import (
	"encoding/json"
	"net/http"

	"go-practice-todo/models"
	"go-practice-todo/services"
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

	tokenUser := &models.User{
		LoginID:  user.LoginID,
		Name:     user.Username,
		Password: user.Password,
	}

	token, err := h.service.GenerateToken(tokenUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	WriteJSONResponse(w, map[string]string{"token": token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		LoginID  string `json:"login_id"`
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
	WriteJSONResponse(w, map[string]string{"token": token})
}
