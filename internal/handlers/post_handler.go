package handlers

import (
	"encoding/json"
	"microblogging-platform/internal/models"
	"microblogging-platform/pkg/logger"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
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

	userID := r.Context().Value("user_id").(uint)
	post.UserID = userID

	if err := h.db.Create(&post).Error; err != nil {
		h.logger.Error("Ошибка при создании поста", err)
		http.Error(w, "Ошибка при создании поста", http.StatusInternalServerError)
		return
	}

	// Создаем новую структуру для ответа
	response := struct {
		ID      uint   `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		ID:      post.ID,
		Title:   post.Content, // Предполагаю, что заголовок хранится в поле Content
		Content: post.Content,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error("Неверный ID поста", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := h.db.First(&post, postID).Error; err != nil {
		h.logger.Error("Пост не найден", err)
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error("Неверный ID поста", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := h.db.First(&post, postID).Error; err != nil {
		h.logger.Error("Пост не найден", err)
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(float64)
	if post.UserID != uint(userID) {
		http.Error(w, "У вас нет прав на редактирование этого поста", http.StatusForbidden)
		return
	}

	var updatedPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		h.logger.Error("Ошибка при декодировании запроса", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	post.Content = updatedPost.Content
	post.ImageURL = updatedPost.ImageURL

	if err := h.db.Save(&post).Error; err != nil {
		h.logger.Error("Ошибка при обновлении поста", err)
		http.Error(w, "Ошибка при обновлении поста", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error("Неверный ID поста", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := h.db.First(&post, postID).Error; err != nil {
		h.logger.Error("Пост не найден", err)
		http.Error(w, "Пост не найден", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(float64)
	if post.UserID != uint(userID) {
		http.Error(w, "У вас нет прав на удаление этого поста", http.StatusForbidden)
		return
	}

	if err := h.db.Delete(&post).Error; err != nil {
		h.logger.Error("Ошибка при удалении поста", err)
		http.Error(w, "Ошибка при удалении поста", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
