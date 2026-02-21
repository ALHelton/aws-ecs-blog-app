package models

import (
	"time"

	"github.com/guregu/null/v6"
)

type BlogPost struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

type Comment struct {
	Id         int64     `json:"id"`
	BlogPostId int64     `json:"blog_post_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  null.Time `json:"deleted_at"`
}
