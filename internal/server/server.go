package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"aws-ecs-blog-app/internal/database"
	"aws-ecs-blog-app/internal/models"

	"gorm.io/gorm"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   *gorm.DB
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db, err := database.NewConnection()
	if err != nil {
		panic(fmt.Sprintf("could not connect to database: %s", err))
	}

	// auto migrate models
	db.AutoMigrate(&models.BlogPost{}, &models.Comment{})

	s := &Server{
		port: port,
		db:   db,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}