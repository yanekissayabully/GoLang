package main

import (
    "fmt"
    "log"
    "net/http"
    "task-api/internal/handlers"
    "task-api/internal/middleware"
    "task-api/internal/models"
    "task-api/internal/storage"
)

func main() {
    store := storage.NewTaskStore()
    handler := handlers.NewTaskHandler(store)
    
    // Создаем маршруты
    mux := http.NewServeMux()
    
    // Регистрируем обработчики
    mux.HandleFunc("GET /tasks", handler.GetTask)
    mux.HandleFunc("POST /tasks", handler.CreateTask)
    mux.HandleFunc("PATCH /tasks", handler.UpdateTask)
    
    // Обертываем в middleware
    stack := middleware.APIKeyMiddleware(middleware.LoggingMiddleware(mux))
    
    // Добавляем несколько тестовых задач (используем models.Task)
    store.Create(models.Task{Title: "Write unit tests", Done: false})
    store.Create(models.Task{Title: "Deploy service", Done: true})
    
    // Запускаем сервер
    port := ":8080"
    fmt.Printf("Server starting on http://localhost%s\n", port)
    fmt.Println("Use API Key: secret12345")
    fmt.Println("Available endpoints:")
    fmt.Println("  GET    /tasks")
    fmt.Println("  GET    /tasks?id=1")
    fmt.Println("  POST   /tasks")
    fmt.Println("  PATCH  /tasks?id=1")
    
    log.Fatal(http.ListenAndServe(port, stack))
}