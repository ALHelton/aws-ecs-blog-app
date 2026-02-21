package server

import (
	"testing"

	"aws-ecs-blog-app/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&models.BlogPost{}, &models.Comment{})
	return db
}