package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"aws-ecs-blog-app/internal/models"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type BlogPostHandler struct {
	db *gorm.DB
}

func NewBlogPostHandler(db *gorm.DB) *BlogPostHandler {
	return &BlogPostHandler{db: db}
}

func (h *BlogPostHandler) GetBlogPosts(w http.ResponseWriter, r *http.Request) {
	var posts []models.BlogPost
	if err := h.db.Find(&posts).Error; err != nil {
		http.Error(w, "error fetching posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *BlogPostHandler) GetBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	var post models.BlogPost
	if err := h.db.First(&post, id).Error; err == gorm.ErrRecordNotFound {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error fetching post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *BlogPostHandler) CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	var post models.BlogPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.db.Create(&post).Error; err != nil {
		http.Error(w, "error creating post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *BlogPostHandler) UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	var post models.BlogPost
	if err := h.db.First(&post, id).Error; err == gorm.ErrRecordNotFound {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error fetching post", http.StatusInternalServerError)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.db.Save(&post).Error; err != nil {
		http.Error(w, "error updating post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *BlogPostHandler) DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid post id", http.StatusBadRequest)
		return
	}

	if err := h.db.Delete(&models.BlogPost{}, id).Error; err == gorm.ErrRecordNotFound {
		http.Error(w, "post not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "error deleting post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}