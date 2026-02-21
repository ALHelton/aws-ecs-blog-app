package handlers

import (
    "sync"

    "aws-ecs-blog-app/internal/models"
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