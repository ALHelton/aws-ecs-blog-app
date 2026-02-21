package server

import (
	"aws-ecs-blog-app/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	healthHandler := handlers.NewHealthHandler()
	r.Get("/", healthHandler.HelloWorld)

	blogPostHandler := handlers.NewBlogPostHandler(s.db)
	r.Get("/blog_posts", blogPostHandler.GetBlogPosts)
	r.Get("/blog_posts/{id}", blogPostHandler.GetBlogPost)
	r.Post("/blog_posts", blogPostHandler.CreateBlogPost)
	r.Put("/blog_posts/{id}", blogPostHandler.UpdateBlogPost)
	r.Delete("/blog_posts/{id}", blogPostHandler.DeleteBlogPost)

	commentHandler := handlers.NewCommentHandler(s.db)
	r.Get("/comments", commentHandler.GetComments)
	r.Get("/comments/{id}", commentHandler.GetComment)
	r.Post("/comments", commentHandler.CreateComment)
	r.Put("/comments/{id}", commentHandler.UpdateComment)
	r.Delete("/comments/{id}", commentHandler.DeleteComment)

	return r
}