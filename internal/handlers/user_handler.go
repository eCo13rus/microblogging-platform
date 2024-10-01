package handlers

import (
	"encoding/json"
	"fmt"
	"microblogging-platform/internal/models"
	"microblogging-platform/pkg/logger"
	"microblogging-platform/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserHandler struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewUserHandler(db *gorm.DB, logger logger.Logger) *UserHandler {
	return &UserHandler{db: db, logger: logger}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := h.db.Find(&users).Error; err != nil {
		h.logger.Error("Ошибка при получении пользователей", err)
		http.Error(w, "Ошибка при получении пользователей", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Получен запрос на получение текущего пользователя")

	userID := r.Context().Value("user_id")

	if userID == nil {
		h.logger.Error("Пользователь не аутентифицирован", nil)
		http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		h.logger.Error("Пользователь не найден", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	h.logger.Info("Пользователь успешно найден и отправлен")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	if userID == "me" {
		// Получаем ID текущего пользователя из контекста
		currentUserID := r.Context().Value("user_id")
		if currentUserID == nil {
			h.logger.Error("Пользователь не аутентифицирован", nil)
			http.Error(w, "Пользователь не аутентифицирован", http.StatusUnauthorized)
			return
		}
		userID = fmt.Sprintf("%v", currentUserID)
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		h.logger.Error("Пользователь не найден", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error("Неверный ID пользователя", err)
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		h.logger.Error("Пользователь не найден", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	currentUserID := r.Context().Value("user_id").(float64)
	if user.ID != uint(currentUserID) {
		http.Error(w, "У вас нет прав на редактирование этого пользователя", http.StatusForbidden)
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		h.logger.Error("Ошибка при декодировании запроса", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateUser(&updatedUser); err != nil {
		h.logger.Error("Ошибка валидации данных пользователя", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Username = updatedUser.Username
	user.Email = updatedUser.Email

	if err := h.db.Save(&user).Error; err != nil {
		h.logger.Error("Ошибка при обновлении пользователя", err)
		http.Error(w, "Ошибка при обновлении пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error("Неверный ID пользователя", err)
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		h.logger.Error("Пользователь не найден", err)
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	currentUserID := r.Context().Value("user_id").(float64)
	if user.ID != uint(currentUserID) {
		http.Error(w, "У вас нет прав на удаление этого пользователя", http.StatusForbidden)
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		h.logger.Error("Ошибка при удалении пользователя", err)
		http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
