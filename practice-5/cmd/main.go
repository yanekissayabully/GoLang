package main

import (
    "net/http"
    "practice5/internal/db"
    "practice5/internal/handlers"
    "practice5/internal/repository"
)

func main() {
    database, _ := db.Connect()
    defer database.Close()  
    repo := repository.NewRepository(database)
    repo.SeedDatabase()
    userHandler := handlers.NewUserHandler(repo)
    http.HandleFunc("/users", userHandler.GetUsersHandler)
    http.HandleFunc("/common-friends", userHandler.GetCommonFriendsHandler)
    http.ListenAndServe(":8080", nil)
}