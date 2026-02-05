package storage

import (
    "sync"
    "task-api/internal/models"
)

type TaskStore struct {
    mu     sync.RWMutex
    tasks  map[int]models.Task
    nextID int
}

func NewTaskStore() *TaskStore {
    return &TaskStore{
        tasks:  make(map[int]models.Task),
        nextID: 1,
    }
}

func (s *TaskStore) Create(task models.Task) models.Task {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    task.ID = s.nextID
    s.tasks[task.ID] = task
    s.nextID++
    return task
}

func (s *TaskStore) GetByID(id int) (models.Task, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    task, exists := s.tasks[id]
    return task, exists
}

func (s *TaskStore) GetAll() []models.Task {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    tasks := make([]models.Task, 0, len(s.tasks))
    for _, task := range s.tasks {
        tasks = append(tasks, task)
    }
    return tasks
}

func (s *TaskStore) Update(id int, done bool) (models.Task, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    task, exists := s.tasks[id]
    if !exists {
        return models.Task{}, false
    }
    
    task.Done = done
    s.tasks[id] = task
    return task, true
}