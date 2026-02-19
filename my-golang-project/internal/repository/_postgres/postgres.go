package _postgres

import (
    "context"
    "fmt"
    "log" 
    "my-golang-project/pkg/modules"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres" 
    _ "github.com/golang-migrate/migrate/v4/source/file"       
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" 
)

//Dialect - наша "обертка" над подключением к БД
type Dialect struct {
    DB *sqlx.DB
}

//NewPGXDialect создает новое подключение и применяет миграции
func NewPGXDialect(ctx context.Context, cfg *modules.PostgreSQL) *Dialect {
    //Формируем строку подключения для sqlx.Connect
    //dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    //cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
        cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

    db, err := sqlx.Connect("postgres", dsn)
    if err != nil {
        log.Fatalf("Ошибка подключения к БД: %v", err) //log.Fatal выведет и завершит программу
    }

    err = db.Ping()
    if err != nil {
        log.Fatalf("Ошибка проверки подключения (Ping): %v", err)
    }

    log.Println("Успешно подключились к БД!")

    //Запускаем миграции
    AutoMigrate(cfg) //cfg передаем, т.к. там есть все данные

    return &Dialect{DB: db}
}

//AutoMigrate применяет миграции из папки database/migrations
func AutoMigrate(cfg *modules.PostgreSQL) {
    sourceURL := "file://database/migrations" //Путь до папки с миграциями
    databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
        cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

    m, err := migrate.New(sourceURL, databaseURL)
    if err != nil {
        log.Fatalf("Ошибка создания объекта миграции: %v", err)
    }

    //Применяем все доступные миграции вверх
    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Ошибка применения миграций: %v", err)
    }

    if err == migrate.ErrNoChange {
        log.Println("Миграции не требуется, БД актуальна")
    } else {
        log.Println("Миграции успешно применены!")
    }
}