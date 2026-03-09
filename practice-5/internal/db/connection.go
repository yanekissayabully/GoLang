package db

import (
    "database/sql"
    _ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
    connStr := "host=localhost port=5432 user=postgres password=admin dbname=practice5 sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    err = db.Ping()
    if err != nil {
        return nil, err
    }

    return db, nil
}