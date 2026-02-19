package modules

import (
    "database/sql"
    "time"
)

type User struct {
    ID        int          `db:"id"`
    Name      string       `db:"name"`
    Email     string       `db:"email"`
    Age       *int         `db:"age"`        // Указатель для NULL
    CreatedAt time.Time    `db:"created_at"`
    DeletedAt sql.NullTime `db:"deleted_at"` // sql.NullTime для NULL timestamp
}

// IsDeleted - удобный метод для проверки, удален ли пользователь
func (u *User) IsDeleted() bool {
    return u.DeletedAt.Valid
}