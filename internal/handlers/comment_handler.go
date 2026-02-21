package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"aws-ecs-blog-app/internal/models"

	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
    mu    sync.Mutex
    comments []models.Comment
}

func NewCommentHandler() *CommentHandler {
    return &CommentHandler{
        comments: []models.Comment{},
    }
}

func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	blogPostId, err := strconv.ParseInt(r.URL.Query().Get("blog_post_id"), 10, 64)
    if err != nil {
        http.Error(w, "invalid blog post id", http.StatusBadRequest)
        return
    }

    h.mu.Lock()
    defer h.mu.Unlock()

    var postComments []models.Comment
    for _, c := range h.comments {
        if c.BlogPostId == blogPostId {
            postComments = append(postComments, c)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(postComments)
}

func (h *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for _, c := range h.comments {
		if c.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "comment not found", http.StatusNotFound)
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

	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()

	h.mu.Lock()
	defer h.mu.Unlock()

	h.comments = append(h.comments, comment)

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

	var updated models.Comment
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, c := range h.comments {
		if c.Id == id {
			h.comments[i].Content = updated.Content
			h.comments[i].UpdatedAt = time.Now()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(h.comments[i])
			return
		}
	}

	http.Error(w, "comment not found", http.StatusNotFound)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid comment id", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for i, c := range h.comments {
		if c.Id == id {
			h.comments = append(h.comments[:i], h.comments[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "comment not found", http.StatusNotFound)
}