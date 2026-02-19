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
        //генерируем Request ID на основе timestamp
        requestID := strconv.FormatInt(time.Now().UnixNano(), 10)
        
        //добавляем в хедеры
        w.Header().Set("X-Request-ID", requestID)
        
        //добавляем в контекст запроса
        ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}