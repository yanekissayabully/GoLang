package http

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "fmt"
    "my-golang-project/internal/usecase"
)

type UserHandler struct {
    usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
    return &UserHandler{usecase: uc}
}

// Request/Response структуры
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

// GetUsers - GET /users (только активные)
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

// GetUserByID - GET /users/{id} (даже удаленных)
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    id, err := extractIDFromPath(r.URL.Path)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
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

// CreateUser - POST /users
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

// UpdateUser - PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    id, err := extractIDFromPath(r.URL.Path)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
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

// DeleteUser - DELETE /users/{id} (мягкое удаление)
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
    id, err := extractIDFromPath(r.URL.Path)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    err = h.usecase.DeleteUser(id)
    if err != nil {
        if strings.Contains(err.Error(), "не найден") || strings.Contains(err.Error(), "уже удален") {
            w.WriteHeader(http.StatusNotFound)
        } else {
            w.WriteHeader(http.StatusInternalServerError)
        }
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "soft deleted"})
}

// ========== НОВЫЕ ЭНДПОИНТЫ ==========

// GetDeletedUsers - GET /users/deleted (получить удаленных)
func (h *UserHandler) GetDeletedUsers(w http.ResponseWriter, r *http.Request) {
    users, err := h.usecase.GetDeletedUsers()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(users)
}

// RestoreUser - POST /users/{id}/restore (восстановить удаленного)
func (h *UserHandler) RestoreUser(w http.ResponseWriter, r *http.Request) {
    id, err := extractIDFromPath(r.URL.Path)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    err = h.usecase.RestoreUser(id)
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
    json.NewEncoder(w).Encode(map[string]string{"status": "restored"})
}

// HardDeleteUser - DELETE /users/{id}/hard (ПОЛНОЕ удаление)
func (h *UserHandler) HardDeleteUser(w http.ResponseWriter, r *http.Request) {
    id, err := extractIDFromPath(r.URL.Path)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
        return
    }

    err = h.usecase.HardDeleteUser(id)
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
    json.NewEncoder(w).Encode(map[string]string{"status": "permanently deleted"})
}

// Вспомогательная функция для извлечения ID из пути
func extractIDFromPath(path string) (int, error) {
    pathParts := strings.Split(path, "/")
    if len(pathParts) < 3 {
        return 0, fmt.Errorf("неверный путь")
    }
    
    // Ищем ID (может быть на разных позициях в зависимости от пути)
    // /users/1 -> id на позиции 2
    // /users/1/restore -> id на позиции 2
    // /users/deleted -> тут ID нет, но этот эндпоинт обрабатывается отдельно
    var idStr string
    if strings.Contains(path, "/restore") || strings.Contains(path, "/hard") {
        idStr = pathParts[2] // /users/1/restore -> 1 на позиции 2
    } else {
        idStr = pathParts[2] // /users/1 -> 1 на позиции 2
    }
    
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return 0, fmt.Errorf("некорректный ID")
    }
    return id, nil
}