package repository

import (
    "my-golang-project/internal/repository/_postgres" // Добавляем этот импорт
    "my-golang-project/internal/repository/_postgres/users"
    "my-golang-project/pkg/modules"
)

type UserRepository interface {
    GetUsers() ([]modules.User, error)
    GetUserByID(id int) (*modules.User, error)
    CreateUser(name, email string, age *int) (int, error) // Обновляем сигнатуру
    UpdateUser(id int, name, email string, age *int) error // Обновляем сигнатуру
    DeleteUser(id int) error
}

type Repositories struct {
    User UserRepository
}

// NewRepositories теперь принимает *_postgres.Dialect
func NewRepositories(db *_postgres.Dialect) *Repositories {
    return &Repositories{
        // Передаем db.DB в конструктор репозитория пользователей
        User: users.NewUserRepository(db.DB),
    }
}