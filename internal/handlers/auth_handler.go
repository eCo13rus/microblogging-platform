package handlers

import (
	"encoding/json"
	"microblogging-platform/internal/config"
	"microblogging-platform/internal/models"
	"microblogging-platform/pkg/logger"
	"microblogging-platform/pkg/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db     *gorm.DB
	logger logger.Logger
	config *config.Config
}

func NewAuthHandler(db *gorm.DB, logger logger.Logger) *AuthHandler {
	return &AuthHandler{db: db, logger: logger}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error("Ошибка при декодировании запроса регистрации", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Валидация данных пользователя
	if err := utils.ValidateUser(&user); err != nil {
		h.logger.Error("Ошибка валидации данных пользователя", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("Ошибка при хешировании пароля", err)
		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Создание пользователя в базе данных
	if err := h.db.Create(&user).Error; err != nil {
		h.logger.Error("Ошибка при создании пользователя", err)
		http.Error(w, "Ошибка при регистрации пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пользователь успешно зарегистрирован"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		h.logger.Error("Ошибка при декодировании запроса авторизации", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		h.logger.Error("Пользователь не найден", err)
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		h.logger.Error("Неверный пароль", err)
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWTToken(user.ID)
	if err != nil {
		h.logger.Error("Ошибка при генерации JWT токена", err)
		http.Error(w, "Ошибка при авторизации", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
