package handlers

import (
	"encoding/json"
	"gorm.io/gorm"
	"microblogging-platform/internal/models"
	"microblogging-platform/pkg/logger"
	"net/http"
)

type PostHandler struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPostHandler(db *gorm.DB, logger logger.Logger) *PostHandler {
	return &PostHandler{db: db, logger: logger}
}

func (h *PostHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.Post
	if err := h.db.Find(&posts).Error; err != nil {
		h.logger.Error("Ошибка при получении постов", err)
		http.Error(w, "Ошибка при получении постов", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		h.logger.Error("Ошибка при декодировании запроса", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Получаем ID первого пользователя из базы данных
	var user models.User
	if err := h.db.First(&user).Error; err != nil {
		h.logger.Error("Ошибка при получении пользователя", err)
		http.Error(w, "Ошибка при создании поста", http.StatusInternalServerError)
		return
	}
	post.UserID = user.ID

	if err := h.db.Create(&post).Error; err != nil {
		h.logger.Error("Ошибка при создании поста", err)
		http.Error(w, "Ошибка при создании поста", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}
