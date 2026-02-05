package main

import (
    "fmt"
    "log"
    "net/http"
    "task-api/internal/external"
    "task-api/internal/handlers"
    "task-api/internal/middleware"
    "task-api/internal/models"
    "task-api/internal/storage"
)

func main() {
    store := storage.NewTaskStore()
    apiClient := external.NewAPIClient()
    handler := handlers.NewTaskHandler(store, apiClient)
    
    store.Create(models.Task{Title: "Write unit tests", Done: false})
    store.Create(models.Task{Title: "Deploy service", Done: true})
    store.Create(models.Task{Title: "Learn Go", Done: false})
    
    mux := http.NewServeMux()
    
    mux.HandleFunc("GET /tasks", handler.GetTask)
    mux.HandleFunc("POST /tasks", handler.CreateTask)
    mux.HandleFunc("PATCH /tasks", handler.UpdateTask)
    mux.HandleFunc("DELETE /tasks", handler.DeleteTask)
    
    mux.HandleFunc("GET /external/todos", handler.GetExternalTodos)
    mux.HandleFunc("POST /external/posts", handler.CreateExternalPost)
    
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"status": "ok", "tasks": %d}`, store.Count())
    })
    
    stack := middleware.APIKeyMiddleware(
        middleware.RequestIDMiddleware(
            middleware.LoggingMiddleware(mux),
        ),
    )
    
    port := ":8080"
    fmt.Printf("Server starting on http://localhost%s\n", port)
    fmt.Println("Use API Key: secret12345")
    fmt.Println("\nAvailable endpoints:")
    fmt.Println("  GET    /tasks                    - Get all tasks")
    fmt.Println("  GET    /tasks?id=1              - Get task by ID")
    fmt.Println("  GET    /tasks?done=true         - Get tasks filtered by status")
    fmt.Println("  POST   /tasks                   - Create new task")
    fmt.Println("  PATCH  /tasks?id=1              - Update task status")
    fmt.Println("  DELETE /tasks?id=1              - Delete task")
    fmt.Println("  GET    /external/todos          - Get todos from external API")
    fmt.Println("  POST   /external/posts          - Create post on external API")
    fmt.Println("  GET    /health                  - Health check")
    
    log.Fatal(http.ListenAndServe(port, stack))
}