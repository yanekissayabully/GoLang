package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "task-api/internal/models"
    "task-api/internal/storage"
)

type TaskHandler struct {
    store *storage.TaskStore
}

func NewTaskHandler(store *storage.TaskStore) *TaskHandler {
    return &TaskHandler{store: store}
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        // Если нет id, возвращаем все задачи
        tasks := h.store.GetAll()
        json.NewEncoder(w).Encode(tasks)
        return
    }
    
    id, err := strconv.Atoi(idStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid id"})
        return
    }
    
    task, exists := h.store.GetByID(id)
    if !exists {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "task not found"})
        return
    }
    
    json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    var req struct {
        Title string `json:"title"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid request body"})
        return
    }
    
    req.Title = strings.TrimSpace(req.Title)
    if req.Title == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid title"})
        return
    }
    
    task := models.Task{
        Title: req.Title,
        Done:  false,
    }
    
    createdTask := h.store.Create(task)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdTask)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "id is required"})
        return
    }
    
    id, err := strconv.Atoi(idStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid id"})
        return
    }
    
    var req struct {
        Done bool `json:"done"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid request body"})
        return
    }
    
    updatedTask, exists := h.store.Update(id, req.Done)
    if !exists {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(models.ErrorResponse{Error: "task not found"})
        return
    }
    
    json.NewEncoder(w).Encode(updatedTask)
}