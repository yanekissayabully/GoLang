package users

import (
    "database/sql"
    "errors"
    "fmt"


    "my-golang-project/pkg/modules"

    "github.com/jmoiron/sqlx"
)

type UserRepositoryPostgres struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepositoryPostgres {
    return &UserRepositoryPostgres{db: db}
}

// GetUsers возвращает только НЕУДАЛЕННЫХ пользователей (deleted_at IS NULL)
func (r *UserRepositoryPostgres) GetUsers() ([]modules.User, error) {
    var users []modules.User
    query := `
        SELECT id, name, email, age, created_at, deleted_at 
        FROM users 
        WHERE deleted_at IS NULL 
        ORDER BY id
    `
    err := r.db.Select(&users, query)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения всех пользователей: %w", err)
    }
    return users, nil
}

// GetUserByID возвращает пользователя по ID (даже если он удален)
// Возвращает ошибку, если пользователь не найден
func (r *UserRepositoryPostgres) GetUserByID(id int) (*modules.User, error) {
    var user modules.User
    query := `
        SELECT id, name, email, age, created_at, deleted_at 
        FROM users 
        WHERE id = $1
    `
    err := r.db.Get(&user, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("пользователь с ID %d не найден", id)
        }
        return nil, fmt.Errorf("ошибка получения пользователя по ID %d: %w", id, err)
    }
    return &user, nil
}

// GetActiveUserByID возвращает только НЕУДАЛЕННОГО пользователя
func (r *UserRepositoryPostgres) GetActiveUserByID(id int) (*modules.User, error) {
    var user modules.User
    query := `
        SELECT id, name, email, age, created_at, deleted_at 
        FROM users 
        WHERE id = $1 AND deleted_at IS NULL
    `
    err := r.db.Get(&user, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("активный пользователь с ID %d не найден", id)
        }
        return nil, fmt.Errorf("ошибка получения пользователя по ID %d: %w", id, err)
    }
    return &user, nil
}

// CreateUser создает нового пользователя
func (r *UserRepositoryPostgres) CreateUser(name, email string, age *int) (int, error) {
    var id int
    query := `
        INSERT INTO users (name, email, age) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `
    err := r.db.QueryRow(query, name, email, age).Scan(&id)
    if err != nil {
        return 0, fmt.Errorf("ошибка создания пользователя: %w", err)
    }
    return id, nil
}

// UpdateUser обновляет данные пользователя (только если он не удален)
func (r *UserRepositoryPostgres) UpdateUser(id int, name, email string, age *int) error {
    // Сначала проверим, существует ли пользователь и не удален ли он
    user, err := r.GetActiveUserByID(id)
    if err != nil {
        return err // пользователь не найден или удален
    }
    _ = user // просто чтобы подавить предупреждение, мы уже проверили существование

    query := `
        UPDATE users 
        SET name = $1, email = $2, age = $3 
        WHERE id = $4 AND deleted_at IS NULL
    `
    result, err := r.db.Exec(query, name, email, age, id)
    if err != nil {
        return fmt.Errorf("ошибка обновления пользователя ID %d: %w", id, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("пользователь с ID %d не найден или уже удален", id)
    }

    return nil
}

// DeleteUser - МЯГКОЕ удаление (ставим deleted_at)
func (r *UserRepositoryPostgres) DeleteUser(id int) error {
    // Проверяем, существует ли пользователь
    user, err := r.GetUserByID(id)
    if err != nil {
        return err // пользователь не найден
    }

    // Если уже удален, возвращаем ошибку
    if user.IsDeleted() {
        return fmt.Errorf("пользователь с ID %d уже удален", id)
    }

    query := `
        UPDATE users 
        SET deleted_at = CURRENT_TIMESTAMP 
        WHERE id = $1 AND deleted_at IS NULL
    `
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("ошибка удаления пользователя ID %d: %w", id, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("пользователь с ID %d не найден", id)
    }

    return nil
}

// HardDeleteUser - ПОЛНОЕ удаление из БД (на всякий случай, для админских функций)
func (r *UserRepositoryPostgres) HardDeleteUser(id int) error {
    query := "DELETE FROM users WHERE id = $1"
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("ошибка полного удаления пользователя ID %d: %w", id, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("пользователь с ID %d не найден", id)
    }

    return nil
}

// GetDeletedUsers - получить всех удаленных пользователей
func (r *UserRepositoryPostgres) GetDeletedUsers() ([]modules.User, error) {
    var users []modules.User
    query := `
        SELECT id, name, email, age, created_at, deleted_at 
        FROM users 
        WHERE deleted_at IS NOT NULL 
        ORDER BY deleted_at DESC
    `
    err := r.db.Select(&users, query)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения удаленных пользователей: %w", err)
    }
    return users, nil
}

// RestoreUser - восстановить удаленного пользователя
func (r *UserRepositoryPostgres) RestoreUser(id int) error {
    query := `
        UPDATE users 
        SET deleted_at = NULL 
        WHERE id = $1 AND deleted_at IS NOT NULL
    `
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("ошибка восстановления пользователя ID %d: %w", id, err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("пользователь с ID %d не найден или не был удален", id)
    }

    return nil
}