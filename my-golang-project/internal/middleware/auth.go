package middleware

import (
    "net/http"
)

const APIKey = "my-secret-key-123"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Пропускаем healthcheck без проверки ключа
        if r.URL.Path == "/health" {
            next.ServeHTTP(w, r)
            return
        }

        key := r.Header.Get("X-API-KEY")
        if key != APIKey {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error":"неавторизован"}`))
            return
        }
        next.ServeHTTP(w, r)
    })
}