package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    _ "github.com/lib/pq"
)

type Movie struct {
    ID          int     `json:"id"`
    Title       string  `json:"title"`
    Genre       string  `json:"genre"`
    Rating      float64 `json:"rating"`
    Description string  `json:"description"`
}

type App struct {
    DB *sql.DB
}

func main() {
    // Ожидание готовности базы данных
    time.Sleep(3 * time.Second)
    
    // Подключение к базе данных
    connStr := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        getEnv("DB_HOST", "db"),
        getEnv("DB_PORT", "5432"),
        getEnv("DB_USER", "postgres"),
        getEnv("DB_PASSWORD", "password"),
        getEnv("DB_NAME", "moviesdb"),
    )

    var db *sql.DB
    var err error
    
    for i := 0; i < 10; i++ {
        db, err = sql.Open("postgres", connStr)
        if err != nil {
            log.Printf("Failed to open DB (attempt %d): %v", i+1, err)
            time.Sleep(2 * time.Second)
            continue
        }
        
        err = db.Ping()
        if err == nil {
            break
        }
        
        log.Printf("Failed to ping DB (attempt %d): %v", i+1, err)
        db.Close()
        time.Sleep(2 * time.Second)
    }
    
    if err != nil {
        log.Fatal("Failed to connect to database after multiple attempts:", err)
    }
    
    defer db.Close()
    log.Println("Database connected successfully!")

    app := &App{DB: db}

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Создание таблицы
    err = app.createTable()
    if err != nil {
        log.Printf("Error creating table: %v", err)
    }

    // Настройка маршрутов
    http.HandleFunc("/movies", app.handleMovies)
    http.HandleFunc("/movies/", app.handleMovieByID)

    server := &http.Server{
        Addr:    ":8000",
        Handler: nil,
    }

    go func() {
        log.Println("Starting the Server on port 8000...")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Could not listen on :8000: %v\n", err)
        }
    }()

    // Ожидание сигнала завершения
    sig := <-sigChan
    log.Printf("Received signal: %s", sig)
    log.Println("Shutting down gracefully...")
    
    // Закрытие соединения с БД
    if err := db.Close(); err != nil {
        log.Printf("Error closing database connection: %v", err)
    }
    log.Println("Database connection closed")
    log.Println("Server stopped")
}

func (app *App) createTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS movies (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        genre VARCHAR(100) NOT NULL,
        rating FLOAT NOT NULL,
        description TEXT
    )`
    
    _, err := app.DB.Exec(query)
    return err
}

func (app *App) handleMovies(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        app.getMovies(w, r)
    case http.MethodPost:
        app.createMovie(w, r)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
    }
}

func (app *App) handleMovieByID(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/movies/"):]
    
    switch r.Method {
    case http.MethodGet:
        app.getMovie(w, r, id)
    case http.MethodPut:
        app.updateMovie(w, r, id)
    case http.MethodDelete:
        app.deleteMovie(w, r, id)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
    }
}

func (app *App) getMovies(w http.ResponseWriter, r *http.Request) {
    rows, err := app.DB.Query("SELECT id, title, genre, rating, description FROM movies")
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
    defer rows.Close()

    movies := []Movie{}
    for rows.Next() {
        var m Movie
        err := rows.Scan(&m.ID, &m.Title, &m.Genre, &m.Rating, &m.Description)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
            return
        }
        movies = append(movies, m)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movies)
}

func (app *App) createMovie(w http.ResponseWriter, r *http.Request) {
    var movie Movie
    err := json.NewDecoder(r.Body).Decode(&movie)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
        return
    }

    var id int
    err = app.DB.QueryRow(
        "INSERT INTO movies (title, genre, rating, description) VALUES ($1, $2, $3, $4) RETURNING id",
        movie.Title, movie.Genre, movie.Rating, movie.Description,
    ).Scan(&id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    movie.ID = id
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(movie)
}

func (app *App) getMovie(w http.ResponseWriter, r *http.Request, id string) {
    var movie Movie
    err := app.DB.QueryRow(
        "SELECT id, title, genre, rating, description FROM movies WHERE id = $1",
        id,
    ).Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Rating, &movie.Description)
    
    if err == sql.ErrNoRows {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
        return
    } else if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movie)
}

func (app *App) updateMovie(w http.ResponseWriter, r *http.Request, id string) {
    var movie Movie
    err := json.NewDecoder(r.Body).Decode(&movie)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
        return
    }

    result, err := app.DB.Exec(
        "UPDATE movies SET title=$1, genre=$2, rating=$3, description=$4 WHERE id=$5",
        movie.Title, movie.Genre, movie.Rating, movie.Description, id,
    )
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
        return
    }

    movie.ID = parseInt(id)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movie)
}

func (app *App) deleteMovie(w http.ResponseWriter, r *http.Request, id string) {
    result, err := app.DB.Exec("DELETE FROM movies WHERE id = $1", id)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        w.WriteHeader(http.StatusNotFound)
        json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}

func parseInt(s string) int {
    var i int
    fmt.Sscanf(s, "%d", &i)
    return i
}