package handlers

import (
    "sync"

    "aws-ecs-blog-app/internal/models"
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