package app

import (
    "context"
    "log"
    "net/http"  // стандартный http пакет
    "time"
	"os"
	"github.com/joho/godotenv"

    userHttp "my-golang-project/internal/delivery/http" // даем алиас userHttp
    "my-golang-project/internal/middleware"
    "my-golang-project/internal/repository"
    "my-golang-project/internal/repository/_postgres"
    // "my-golang-project/internal/repository/_postgres/users"
    "my-golang-project/internal/usecase"
    "my-golang-project/pkg/modules"
)

func Run() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    log.Println("Инициализация конфига БД...")
    dbConfig := initPostgreSQL()

    log.Println("Подключение к БД и применение миграций...")
    postgreDialect := _postgres.NewPGXDialect(ctx, dbConfig)

    // --- Инициализация слоев ---
    log.Println("Инициализация репозиториев...")
    // Создаем конкретный репозиторий пользователей
    // userRepo := users.NewUserRepository(postgreDialect.DB)
    // _ = userRepo
    
    // Создаем обертку Repositories
    repos := repository.NewRepositories(postgreDialect)

    // Инициализация Usecase (передаем интерфейс, но передаем конкретную реализацию)
    userUsecase := usecase.NewUserUsecase(repos.User)

    // Инициализация Handler (используем алиас userHttp)
    userHandler := userHttp.NewUserHandler(userUsecase)

    // --- Настройка маршрутов ---
    mux := http.NewServeMux()

    // Публичный healthcheck (до аутентификации)
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })

    // Группа маршрутов для пользователей (защищенных)
    mux.HandleFunc("GET /users", userHandler.GetUsers)
    mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
    mux.HandleFunc("POST /users", userHandler.CreateUser)
    mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
    mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

    // Применяем мидлвари
    handlerWithMiddleware := middleware.LoggingMiddleware(mux) // сначала логируем
    handlerWithMiddleware = middleware.AuthMiddleware(handlerWithMiddleware) // потом проверяем ключ

    // Запуск сервера
    server := &http.Server{
        Addr:         ":8080",
        Handler:      handlerWithMiddleware,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    log.Println("Сервер запущен на :8080")
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }
}

func initPostgreSQL() *modules.PostgreSQL {
    // Загружаем .env файл
    if err := godotenv.Load(); err != nil {
        log.Println("Файл .env не найден, используем значения по умолчанию")
    }
    
    return &modules.PostgreSQL{
        Host:        getEnv("DB_HOST", "localhost"),
        Port:        getEnv("DB_PORT", "5432"),
        Username:    getEnv("DB_USER", "postgres"),
        Password:    getEnv("DB_PASSWORD", "postgres"),
        DBName:      getEnv("DB_NAME", "mydb"),
        SSLMode:     getEnv("DB_SSLMODE", "disable"),
        ExecTimeout: 5 * time.Second,
    }
}

func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}