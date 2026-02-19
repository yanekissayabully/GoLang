package repository

import (
    "my-golang-project/internal/repository/_postgres"
    "my-golang-project/internal/repository/_postgres/users"
    "my-golang-project/pkg/modules"
)

type UserRepository interface {
    GetUsers() ([]modules.User, error)
    GetUserByID(id int) (*modules.User, error)
    GetActiveUserByID(id int) (*modules.User, error) // новый метод
    CreateUser(name, email string, age *int) (int, error)
    UpdateUser(id int, name, email string, age *int) error
    DeleteUser(id int) error                         // мягкое удаление
    HardDeleteUser(id int) error                      // полное удаление
    GetDeletedUsers() ([]modules.User, error)         // получить удаленных
    RestoreUser(id int) error                          // восстановить
}

type Repositories struct {
    User UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
    return &Repositories{
        User: users.NewUserRepository(db.DB),
    }
}