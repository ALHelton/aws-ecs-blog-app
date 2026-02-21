package handlers

import (
	"aws-ecs-blog-app/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type BlogPostHandler struct {
	mu    sync.Mutex
	posts []models.BlogPost
}

func NewBlogPostHandler() *BlogPostHandler {
	return &BlogPostHandler{
		posts: []models.BlogPost{},
	}
}

func (h *BlogPostHandler) GetBlogPosts(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.posts)
}

func (h *BlogPostHandler) GetBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid blog post id", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, p := range h.posts {
		if p.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "blog post not found", http.StatusNotFound)
}

func (h *BlogPostHandler) CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	var post models.BlogPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	post.Id = int64(len(h.posts) + 1)
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	h.mu.Lock()
	defer h.mu.Unlock()

	h.posts = append(h.posts, post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *BlogPostHandler) UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid blog post id", http.StatusBadRequest)
		return
	}

	var updated models.BlogPost
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, p := range h.posts {
		if p.Id == id {
			h.posts[i].Title = updated.Title
			h.posts[i].Content = updated.Content
			h.posts[i].UpdatedAt = time.Now()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(h.posts[i])
			return
		}
	}

	http.Error(w, "blog post not found", http.StatusNotFound)
}

func (h *BlogPostHandler) DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid blog post id", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, p := range h.posts {
		if p.Id == id {
			h.posts = append(h.posts[:i], h.posts[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "blog post not found", http.StatusNotFound)
}
