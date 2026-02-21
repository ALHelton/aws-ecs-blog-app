package server

import (
	"aws-ecs-blog-app/internal/handlers"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestHelloWorldHandler(t *testing.T) {
	h := handlers.NewHealthHandler()
	server := httptest.NewServer(http.HandlerFunc(h.HelloWorld))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	expected := "{\"message\":\"Hello World\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

func TestGetBlogPosts(t *testing.T) {
	h := handlers.NewBlogPostHandler()
	server := httptest.NewServer(http.HandlerFunc(h.GetBlogPosts))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}

func TestCreateBlogPost(t *testing.T) {
	h := handlers.NewBlogPostHandler()
	server := httptest.NewServer(http.HandlerFunc(h.CreateBlogPost))
	defer server.Close()

	blogPost := map[string]string{"title": "Test Blog Post", "content": "Test Content"}
	body, _ := json.Marshal(blogPost)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status Created; got %v", resp.Status)
	}
}

func TestGetBlogPost(t *testing.T) {
	h := handlers.NewBlogPostHandler()

	createServer := httptest.NewServer(http.HandlerFunc(h.CreateBlogPost))
	blogPost := map[string]string{"title": "Test Blog Post", "content": "Test Content"}
	body, _ := json.Marshal(blogPost)
	http.Post(createServer.URL, "application/json", bytes.NewBuffer(body))
	createServer.Close()

	router := chi.NewRouter()
	router.Get("/blog_posts/{id}", h.GetBlogPost)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/blog_posts/1")
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}

func TestUpdateBlogPost(t *testing.T) {
	h := handlers.NewBlogPostHandler()

	createServer := httptest.NewServer(http.HandlerFunc(h.CreateBlogPost))
	blogPost := map[string]string{"title": "Test Blog Post", "content": "Test Content"}
	body, _ := json.Marshal(blogPost)
	http.Post(createServer.URL, "application/json", bytes.NewBuffer(body))
	createServer.Close()

	router := chi.NewRouter()
	router.Put("/blog_posts/{id}", h.UpdateBlogPost)
	server := httptest.NewServer(router)
	defer server.Close()

	updated := map[string]string{"title": "Updated Title", "content": "Updated Content"}
	body, _ = json.Marshal(updated)

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/blog_posts/1", bytes.NewBuffer(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}

func TestDeleteBlogPost(t *testing.T) {
	h := handlers.NewBlogPostHandler()

	createServer := httptest.NewServer(http.HandlerFunc(h.CreateBlogPost))
	blogPost := map[string]string{"title": "Test Blog Post", "content": "Test Content"}
	body, _ := json.Marshal(blogPost)
	http.Post(createServer.URL, "application/json", bytes.NewBuffer(body))
	createServer.Close()

	router := chi.NewRouter()
	router.Delete("/blog_posts/{id}", h.DeleteBlogPost)
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/blog_posts/1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status NoContent; got %v", resp.Status)
	}
}

func TestGetComments(t *testing.T) {
	h := handlers.NewCommentHandler()
	server := httptest.NewServer(http.HandlerFunc(h.GetComments))
	defer server.Close()

	resp, err := http.Get(server.URL + "?blog_post_id=1")
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}

func TestCreateComment(t *testing.T) {
	h := handlers.NewCommentHandler()
	server := httptest.NewServer(http.HandlerFunc(h.CreateComment))
	defer server.Close()

	comment := map[string]interface{}{"blog_post_id": 1, "content": "Test Comment"}
	body, _ := json.Marshal(comment)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status Created; got %v", resp.Status)
	}
}

func TestGetComment(t *testing.T) {
	h := handlers.NewCommentHandler()

	createServer := httptest.NewServer(http.HandlerFunc(h.CreateComment))
	comment := map[string]interface{}{"blog_post_id": 1, "content": "Test Comment"}
	body, _ := json.Marshal(comment)
	http.Post(createServer.URL, "application/json", bytes.NewBuffer(body))
	createServer.Close()

	router := chi.NewRouter()
	router.Get("/comments/{id}", h.GetComment)
	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/comments/1")
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}

func TestUpdateComment(t *testing.T) {
	h := handlers.NewCommentHandler()

	createServer := httptest.NewServer(http.HandlerFunc(h.CreateComment))
	comment := map[string]interface{}{"blog_post_id": 1, "content": "Test Comment"}
	body, _ := json.Marshal(comment)
	http.Post(createServer.URL, "application/json", bytes.NewBuffer(body))
	createServer.Close()

	router := chi.NewRouter()
	router.Put("/comments/{id}", h.UpdateComment)
	server := httptest.NewServer(router)
	defer server.Close()

	updated := map[string]string{"content": "Updated Comment"}
	body, _ = json.Marshal(updated)

	req, _ := http.NewRequest(http.MethodPut, server.URL+"/comments/1", bytes.NewBuffer(body))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
}

func TestDeleteComment(t *testing.T) {
	h := handlers.NewCommentHandler()

	createServer := httptest.NewServer(http.HandlerFunc(h.CreateComment))
	comment := map[string]interface{}{"blog_post_id": 1, "content": "Test Comment"}
	body, _ := json.Marshal(comment)
	http.Post(createServer.URL, "application/json", bytes.NewBuffer(body))
	createServer.Close()

	router := chi.NewRouter()
	router.Delete("/comments/{id}", h.DeleteComment)
	server := httptest.NewServer(router)
	defer server.Close()

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/comments/1", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status NoContent; got %v", resp.Status)
	}
}