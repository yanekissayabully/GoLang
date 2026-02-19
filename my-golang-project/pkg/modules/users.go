package modules

import "time"

type User struct {
    ID        int        `db:"id"`
    Name      string     `db:"name"`
    Email     string     `db:"email"`      // Новое поле
    Age       *int       `db:"age"`        // Указатель, чтобы можно было NULL (если возраст не указан)
    CreatedAt time.Time  `db:"created_at"` // Новое поле. time.Time умеет работать с timestamp
}