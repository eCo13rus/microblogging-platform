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

type LikeHandler struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewLikeHandler(db *gorm.DB, logger logger.Logger) *LikeHandler {
	return &LikeHandler{db: db, logger: logger}
}

func (h *LikeHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postId"])
	if err != nil {
		h.logger.Error("Неверный ID поста", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(float64)

	like := models.Like{
		PostID: uint(postID),
		UserID: uint(userID),
	}

	// Проверяем, существует ли уже лайк
	var existingLike models.Like
	if err := h.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&existingLike).Error; err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Пост уже лайкнут"})
		return
	}

	if err := h.db.Create(&like).Error; err != nil {
		h.logger.Error("Ошибка при создании лайка", err)
		http.Error(w, "Ошибка при лайке поста", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Пост успешно лайкнут"})
}

func (h *LikeHandler) UnlikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postId"])
	if err != nil {
		h.logger.Error("Неверный ID поста", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(float64)

	if err := h.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.Like{}).Error; err != nil {
		h.logger.Error("Ошибка при удалении лайка", err)
		http.Error(w, "Ошибка при удалении лайка", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Лайк успешно удален"})
}
