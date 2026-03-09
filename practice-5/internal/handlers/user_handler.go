package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "practice5/internal/models"
    "practice5/internal/repository"
    "github.com/google/uuid"
)

type UserHandler struct {
    repo *repository.Repository
}

func NewUserHandler(repo *repository.Repository) *UserHandler {
    return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 {
        page = 1
    }
    
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
    if pageSize < 1 {
        pageSize = 10
    }
    
    filters := models.FilterParams{}
    
    if idStr := r.URL.Query().Get("id"); idStr != "" {
        id, _ := uuid.Parse(idStr)
        filters.ID = &id
    }
    
    if name := r.URL.Query().Get("name"); name != "" {
        filters.Name = &name
    }
    
    if email := r.URL.Query().Get("email"); email != "" {
        filters.Email = &email
    }
    
    if gender := r.URL.Query().Get("gender"); gender != "" {
        filters.Gender = &gender
    }
    
    sort := models.SortParams{
        Field:     r.URL.Query().Get("sortBy"),
        Direction: r.URL.Query().Get("sortDir"),
    }
    
    response, _ := h.repo.GetPaginatedUsers(page, pageSize, filters, sort)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) GetCommonFriendsHandler(w http.ResponseWriter, r *http.Request) {
    user1Str := r.URL.Query().Get("user1")
    user2Str := r.URL.Query().Get("user2")
    
    user1ID, _ := uuid.Parse(user1Str)
    user2ID, _ := uuid.Parse(user2Str)
    
    commonFriends, _ := h.repo.GetCommonFriends(user1ID, user2ID)
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(commonFriends)
}