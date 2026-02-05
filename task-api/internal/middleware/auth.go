package middleware

import (
    "encoding/json"
    "net/http"
    "task-api/internal/models"
)

const ValidAPIKey = "secret12345"

func APIKeyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-KEY")
        if apiKey != ValidAPIKey {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(models.ErrorResponse{
                Error:   "unauthorized",
                Details: "invalid or missing API key",
            })
            return
        }
        next.ServeHTTP(w, r)
    })
}