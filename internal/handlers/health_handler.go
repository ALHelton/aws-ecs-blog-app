package handlers

import (
    "encoding/json"
    "log"
    "net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
    return &HealthHandler{}
}

func (h *HealthHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
    resp := make(map[string]string)
    resp["message"] = "Hello World"

    jsonResp, err := json.Marshal(resp)
    if err != nil {
        log.Fatalf("error handling JSON marshal. Err: %v", err)
    }

    _, _ = w.Write(jsonResp)
}