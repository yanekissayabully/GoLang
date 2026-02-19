package http

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"

    "my-golang-project/internal/usecase"
    // Убираем этот импорт, если он действительно не используется напрямую
    // "my-golang-project/pkg/modules"
)

type UserHandler struct {
    usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
    return &UserHandler{usecase: uc}
}

// --- Request/Response структуры ---
type createUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   *int   `json:"age"`
}

type updateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   *int   `json:"age"`
}

type errorResponse struct {
    Error string `json:"error"`
}

// --- Методы-хендлеры ---

// GetUsers обрабатывает GET /users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.usecase.GetUsers()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(users)
}

// GetUserByID обрабатывает GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    // Достаем ID из URL
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 3 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "неверный путь"})
        return
    }
    idStr := pathParts[2]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "некорректный ID"})
        return
    }

    user, err := h.usecase.GetUserByID(id)
    if err != nil {
        if strings.Contains(err.Error(), "не найден") {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}

// CreateUser обрабатывает POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req createUserRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "неверный формат JSON"})
        return
    }

    id, err := h.usecase.CreateUser(req.Name, req.Email, req.Age)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// UpdateUser обрабатывает PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    // Парсим ID
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 3 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "неверный путь"})
        return
    }
    idStr := pathParts[2]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "некорректный ID"})
        return
    }

    var req updateUserRequest
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "неверный формат JSON"})
        return
    }

    err = h.usecase.UpdateUser(id, req.Name, req.Email, req.Age)
    if err != nil {
        if strings.Contains(err.Error(), "не найден") {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusBadRequest)
        }
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// DeleteUser обрабатывает DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    pathParts := strings.Split(r.URL.Path, "/")
    if len(pathParts) < 3 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "неверный путь"})
        return
    }
    idStr := pathParts[2]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: "некорректный ID"})
        return
    }

    err = h.usecase.DeleteUser(id)
    if err != nil {
        if strings.Contains(err.Error(), "не найден") {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}