package middleware

import (
    "context"
    "net/http"
    "strconv"
    "time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Генерируем Request ID на основе timestamp
        requestID := strconv.FormatInt(time.Now().UnixNano(), 10)
        
        // Добавляем в заголовки ответа
        w.Header().Set("X-Request-ID", requestID)
        
        // Добавляем в контекст запроса
        ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}