package users

import (
    "database/sql"
    "errors"
    "fmt"

    "my-golang-project/pkg/modules"

    "github.com/jmoiron/sqlx"
)

// UserRepositoryPostgres - конкретная реализация UserRepository для PostgreSQL
type UserRepositoryPostgres struct {
    db *sqlx.DB
}

// NewUserRepository - конструктор. Принимает *sqlx.DB, возвращает наш репозиторий
func NewUserRepository(db *sqlx.DB) *UserRepositoryPostgres {
    return &UserRepositoryPostgres{db: db}
}

// GetUsers возвращает всех пользователей
func (r *UserRepositoryPostgres) GetUsers() ([]modules.User, error) {
    var users []modules.User
    query := "SELECT id, name, email, age, created_at FROM users ORDER BY id"
    err := r.db.Select(&users, query) // sqlx умеет сканировать сразу в срез структур
    if err != nil {
        return nil, fmt.Errorf("ошибка получения всех пользователей: %w", err)
    }
    return users, nil
}

// GetUserByID возвращает одного пользователя по ID
func (r *UserRepositoryPostgres) GetUserByID(id int) (*modules.User, error) {
    var user modules.User
    query := "SELECT id, name, email, age, created_at FROM users WHERE id = $1"
    err := r.db.Get(&user, query, id) // sqlx сканирует одну строку в структуру
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // Возвращаем nil и понятную ошибку, если пользователь не найден
            return nil, fmt.Errorf("пользователь с ID %d не найден", id)
        }
        return nil, fmt.Errorf("ошибка получения пользователя по ID %d: %w", id, err)
    }
    return &user, nil
}

// CreateUser создает нового пользователя и возвращает его ID
// Меняем сигнатуру функции: CreateUser(name, email string, age *int) (int, error)
func (r *UserRepositoryPostgres) CreateUser(name, email string, age *int) (int, error) {
    var id int
    query := "INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id"
    // age может быть nil, это нормально для *int
    err := r.db.QueryRow(query, name, email, age).Scan(&id)
    if err != nil {
        // Проверим на нарушение уникальности email (код 23505 в PostgreSQL)
        // Это уже продвинутая обработка, пока можно просто вернуть ошибку
        return 0, fmt.Errorf("ошибка создания пользователя: %w", err)
    }
    return id, nil
}

// UpdateUser обновляет имя пользователя
func (r *UserRepositoryPostgres) UpdateUser(id int, name, email string, age *int)  error {
    query := "UPDATE users SET name = $1, email = $2, age = $3 WHERE id = $4"
    result, err := r.db.Exec(query, name, email, age, id)
    if err != nil {
        return fmt.Errorf("ошибка обновления пользователя ID %d: %w", id, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка получения количества затронутых строк при обновлении ID %d: %w", id, err)
    }

    if rowsAffected == 0 {
        // Пользователь не найден, возвращаем кастомную ошибку
        return fmt.Errorf("пользователь с ID %d не найден для обновления", id)
    }

    return nil
}

// DeleteUser удаляет пользователя
func (r *UserRepositoryPostgres) DeleteUser(id int) error {
    query := "DELETE FROM users WHERE id = $1"
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("ошибка удаления пользователя ID %d: %w", id, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка получения количества затронутых строк при удалении ID %d: %w", id, err)
    }

    if rowsAffected == 0 {
        // Пользователь не найден, возвращаем кастомную ошибку
        return fmt.Errorf("пользователь с ID %d не найден для удаления", id)
    }

    return nil
}