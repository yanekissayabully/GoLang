package middleware

import (
    "net/http"
    "os"
)

//GetAPIKey возвращает API ключ из переменных окружения
func GetAPIKey() string {
    if key := os.Getenv("API_KEY"); key != "" {
        return key
    }
    //Значение по умолчанию, если не задано в .env
    return "my-secret-key-123"
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //скипаем healthcheck без проверки ключа
        if r.URL.Path == "/health" {
            next.ServeHTTP(w, r)
            return
        }

        apiKey := GetAPIKey()
        key := r.Header.Get("X-API-KEY")
        
        if key != apiKey {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error":"неавторизован"}`))
            return
        }
        next.ServeHTTP(w, r)
    })
}