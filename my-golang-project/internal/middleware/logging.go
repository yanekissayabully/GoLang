package middleware

import (
    "log"
    "net/http"
    "time"
)

//LoggingMiddleware логирует каждый запрос
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        //Логируем
        log.Printf("[%s] %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

        //Вызываем следующий обработчик
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        log.Printf("Запрос обработан за %v", duration)
    })
}