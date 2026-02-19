package middleware

import (
    "log"
    "net/http"
    "time"
)

// LoggingMiddleware логирует каждый запрос
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        // Логируем перед обработкой
        log.Printf("[%s] %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)

        // Вызываем следующий обработчик
        next.ServeHTTP(w, r)

        // Можно добавить логирование после, например, статус код, но для этого нужен кастомный ResponseWriter
        // Пока оставим так.
        duration := time.Since(start)
        log.Printf("Запрос обработан за %v", duration)
    })
}