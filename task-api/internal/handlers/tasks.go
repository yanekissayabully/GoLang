package handlers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "task-api/internal/external"
    "task-api/internal/middleware"
    "task-api/internal/models"
    "task-api/internal/storage"
)

type TaskHandler struct {
    store     *storage.TaskStore
    apiClient *external.APIClient
}

func NewTaskHandler(store *storage.TaskStore, apiClient *external.APIClient) *TaskHandler {
    return &TaskHandler{
        store:     store,
        apiClient: apiClient,
    }
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    idStr := r.URL.Query().Get("id")
    
    if idStr != "" {
        id, err := strconv.Atoi(idStr)
        if err != nil || id <= 0 {
            w.WriteHeader(http.StatusBadRequest)
            sendError(w, r, "invalid id", "id must be a positive integer", nil)
            return
        }
        
        task, exists := h.store.GetByID(id)
        if !exists {
            w.WriteHeader(http.StatusNotFound)
            sendError(w, r, "task not found", fmt.Sprintf("task with id %d does not exist", id), nil)
            return
        }
        
        json.NewEncoder(w).Encode(task)
        return
    }
    
    doneFilter := r.URL.Query().Get("done")
    var filterDone *bool
    
    if doneFilter != "" {
        parsedDone, err := strconv.ParseBool(doneFilter)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            sendError(w, r, "invalid filter parameter", "done parameter must be 'true' or 'false'", nil)
            return
        }
        filterDone = &parsedDone
    }
    
    tasks := h.store.GetAllFiltered(filterDone)
    json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    var req struct {
        Title string `json:"title"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        sendError(w, r, "invalid request body", "expected JSON with 'title' field", nil)
        return
    }
    
    req.Title = strings.TrimSpace(req.Title)
    validationErrors := []models.ValidationError{}
    
    if req.Title == "" {
        validationErrors = append(validationErrors, models.ValidationError{
            Field:   "title",
            Message: "title cannot be empty",
        })
    }
    
    if len(req.Title) > 100 {
        validationErrors = append(validationErrors, models.ValidationError{
            Field:   "title",
            Message: "title too long, maximum 100 characters",
        })
    }
    
    if len(validationErrors) > 0 {
        w.WriteHeader(http.StatusBadRequest)
        sendError(w, r, "validation failed", 
            fmt.Sprintf("found %d validation error(s)", len(validationErrors)), 
            validationErrors)
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
        sendError(w, r, "missing parameter", "id query parameter is required", nil)
        return
    }
    
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        w.WriteHeader(http.StatusBadRequest)
        sendError(w, r, "invalid id", "id must be a positive integer", nil)
        return
    }
    
    var req struct {
        Done bool `json:"done"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        sendError(w, r, "invalid request body", "expected JSON with 'done' field (boolean)", nil)
        return
    }
    
    updatedTask, exists := h.store.Update(id, req.Done)
    if !exists {
        w.WriteHeader(http.StatusNotFound)
        sendError(w, r, "task not found", fmt.Sprintf("task with id %d does not exist", id), nil)
        return
    }
    
    json.NewEncoder(w).Encode(updatedTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        w.WriteHeader(http.StatusBadRequest)
        sendError(w, r, "missing parameter", "id query parameter is required", nil)
        return
    }
    
    id, err := strconv.Atoi(idStr)
    if err != nil || id <= 0 {
        w.WriteHeader(http.StatusBadRequest)
        sendError(w, r, "invalid id", "id must be a positive integer", nil)
        return
    }
    
    deleted := h.store.Delete(id)
    if !deleted {
        w.WriteHeader(http.StatusNotFound)
        sendError(w, r, "task not found", fmt.Sprintf("task with id %d does not exist", id), nil)
        return
    }
    
    json.NewEncoder(w).Encode(models.SuccessResponse{
        Message: "task deleted successfully",
        Deleted: true,
        ID:      id,
    })
}

func (h *TaskHandler) GetExternalTodos(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    todos, err := h.apiClient.GetExternalTodos()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        sendError(w, r, "failed to fetch external todos", err.Error(), nil)
        return
    }
    
    if len(todos) > 10 {
        todos = todos[:10]
    }
    
    json.NewEncoder(w).Encode(todos)
}

func (h *TaskHandler) CreateExternalPost(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    var req models.CreatePostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        sendError(w, r, "invalid request body", "expected JSON with title, body, and userId fields", nil)
        return
    }
    
    if req.Title == "" || req.Body == "" || req.UserID <= 0 {
        w.WriteHeader(http.StatusBadRequest)
        sendError(w, r, "validation failed", "title, body, and userId are required", nil)
        return
    }
    
    post, err := h.apiClient.CreateExternalPost(req)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        sendError(w, r, "failed to create external post", err.Error(), nil)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(post)
}

func sendError(w http.ResponseWriter, r *http.Request, errorMsg, details string, validationErrors []models.ValidationError) {
    errorResponse := models.DetailedErrorResponse{
        Error:   errorMsg,
        Details: details,
        Errors:  validationErrors,
    }
    
    if requestID, ok := r.Context().Value(middleware.RequestIDKey).(string); ok {
        errorResponse.RequestID = requestID
    }
    
    json.NewEncoder(w).Encode(errorResponse)
}