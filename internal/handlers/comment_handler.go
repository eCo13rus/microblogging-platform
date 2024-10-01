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

type CommentHandler struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewCommentHandler(db *gorm.DB, logger logger.Logger) *CommentHandler {
	return &CommentHandler{db: db, logger: logger}
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postId"])
	if err != nil {
		h.logger.Error("Неверный ID поста", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	var comments []models.Comment
	if err := h.db.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		h.logger.Error("Ошибка при получении комментариев", err)
		http.Error(w, "Ошибка при получении комментариев", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postId"])
	if err != nil {
		h.logger.Error("Неверный ID поста", err)
		http.Error(w, "Неверный ID поста", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		h.logger.Error("Ошибка при декодировании запроса", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("user_id").(float64)
	comment.UserID = uint(userID)
	comment.PostID = uint(postID)

	if err := h.db.Create(&comment).Error; err != nil {
		h.logger.Error("Ошибка при создании комментария", err)
		http.Error(w, "Ошибка при создании комментария", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error("Неверный ID комментария", err)
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := h.db.First(&comment, commentID).Error; err != nil {
		h.logger.Error("Комментарий не найден", err)
		http.Error(w, "Комментарий не найден", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(float64)
	if comment.UserID != uint(userID) {
		http.Error(w, "У вас нет прав на редактирование этого комментария", http.StatusForbidden)
		return
	}

	var updatedComment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&updatedComment); err != nil {
		h.logger.Error("Ошибка при декодировании запроса", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	comment.Content = updatedComment.Content

	if err := h.db.Save(&comment).Error; err != nil {
		h.logger.Error("Ошибка при обновлении комментария", err)
		http.Error(w, "Ошибка при обновлении комментария", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.logger.Error("Неверный ID комментария", err)
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := h.db.First(&comment, commentID).Error; err != nil {
		h.logger.Error("Комментарий не найден", err)
		http.Error(w, "Комментарий не найден", http.StatusNotFound)
		return
	}

	userID := r.Context().Value("user_id").(float64)
	if comment.UserID != uint(userID) {
		http.Error(w, "У вас нет прав на удаление этого комментария", http.StatusForbidden)
		return
	}

	if err := h.db.Delete(&comment).Error; err != nil {
		h.logger.Error("Ошибка при удалении комментария", err)
		http.Error(w, "Ошибка при удалении комментария", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
