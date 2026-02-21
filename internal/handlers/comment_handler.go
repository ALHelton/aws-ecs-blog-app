package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"aws-ecs-blog-app/internal/models"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type CommentHandler struct {
	db *gorm.DB
}

func NewCommentHandler(db *gorm.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	blogPostId, err := strconv.ParseInt(r.URL.Query().Get("blog_post_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid blog post id", http.StatusBadRequest)
		return
	}

	var comments []models.Comment
	if err := h.db.Where("blog_post_id = ?", blogPostId).Find(&comments).Error; err != nil {
		http.Error(w, "error fetching comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := h.db.First(&comment, id).Error; err == gorm.ErrRecordNotFound {
		http.Error(w, "comment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error fetching comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if comment.BlogPostId == 0 {
		http.Error(w, "blog_post_id is required", http.StatusBadRequest)
		return
	}

	if err := h.db.Create(&comment).Error; err != nil {
		http.Error(w, "error creating comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	if err := h.db.First(&comment, id).Error; err == gorm.ErrRecordNotFound {
		http.Error(w, "comment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error fetching comment", http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.db.Save(&comment).Error; err != nil {
		http.Error(w, "error updating comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comment)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	if err := h.db.Delete(&models.Comment{}, id).Error; err == gorm.ErrRecordNotFound {
		http.Error(w, "comment not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error deleting comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}