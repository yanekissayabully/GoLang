package app

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"
    "github.com/joho/godotenv"

    userHttp "my-golang-project/internal/delivery/http"
    "my-golang-project/internal/middleware"
    "my-golang-project/internal/repository"
    "my-golang-project/internal/repository/_postgres"
    "my-golang-project/internal/usecase"
    "my-golang-project/pkg/modules"
)

func Run() {
    if err := godotenv.Load(); err != nil {
        log.Println("Файл .env не найден, используем переменные окружения или значения по умолчанию")
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    log.Println("Инициализация конфига БД...")
    dbConfig := initPostgreSQL()

    log.Println("Подключение к БД и применение миграций...")
    postgreDialect := _postgres.NewPGXDialect(ctx, dbConfig)

    log.Println("Инициализация репозиториев...")
    repos := repository.NewRepositories(postgreDialect)

    userUsecase := usecase.NewUserUsecase(repos.User)
    userHandler := userHttp.NewUserHandler(userUsecase)

    //Настройка маршрутов
    mux := http.NewServeMux()

    //Публичный healthcheck
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })

    //Основные маршруты
    mux.HandleFunc("GET /users", userHandler.GetUsers)
    mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
    mux.HandleFunc("POST /users", userHandler.CreateUser)
    mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUser)
    mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUser)

    //Маршруты для soft delete
    mux.HandleFunc("GET /users/deleted", userHandler.GetDeletedUsers)
    mux.HandleFunc("POST /users/{id}/restore", userHandler.RestoreUser)
    mux.HandleFunc("DELETE /users/{id}/hard", userHandler.HardDeleteUser)

    //Применяем мидлвари
    handlerWithMiddleware := middleware.LoggingMiddleware(mux)
    handlerWithMiddleware = middleware.AuthMiddleware(handlerWithMiddleware)

    //Получаем порт из .env или используем 8080 по умолчанию
    port := getEnv("SERVER_PORT", "8080")
    serverAddr := ":" + port

    //Таймауты из .env
    readTimeout := getEnvAsInt("SERVER_READ_TIMEOUT", 10)
    writeTimeout := getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)
    idleTimeout := getEnvAsInt("SERVER_IDLE_TIMEOUT", 120)

    //Запуск сервера с graceful shutdown
    server := &http.Server{
        Addr:         serverAddr,
        Handler:      handlerWithMiddleware,
        ReadTimeout:  time.Duration(readTimeout) * time.Second,
        WriteTimeout: time.Duration(writeTimeout) * time.Second,
        IdleTimeout:  time.Duration(idleTimeout) * time.Second,
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Printf("Сервер запущен на %s", serverAddr)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Ошибка запуска сервера: %v", err)
        }
    }()

    <-quit
    log.Println("Получен сигнал завершения, останавливаем сервер...")

    ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancelShutdown()

    if err := server.Shutdown(ctxShutdown); err != nil {
        log.Fatalf("Ошибка при остановке сервера: %v", err)
    }

    if err := postgreDialect.DB.Close(); err != nil {
        log.Printf("Ошибка при закрытии БД: %v", err)
    }

    log.Println("Сервер успешно остановлен")
}

func initPostgreSQL() *modules.PostgreSQL {
    timeout := getEnvAsInt("DB_TIMEOUT", 5)
    
    return &modules.PostgreSQL{
        Host:        getEnv("DB_HOST", "localhost"),
        Port:        getEnv("DB_PORT", "5432"),
        Username:    getEnv("DB_USER", "postgres"),
        Password:    getEnv("DB_PASSWORD", "postgres"),
        DBName:      getEnv("DB_NAME", "mydb"),
        SSLMode:     getEnv("DB_SSLMODE", "disable"),
        ExecTimeout: time.Duration(timeout) * time.Second,
    }
}

//Вспомогательные функции для работы с .env
func getEnv(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    valueStr := getEnv(key, "")
    if value, err := strconv.Atoi(valueStr); err == nil {
        return value
    }
    return defaultValue
}