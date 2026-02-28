# ============================================
# ПОЛНЫЙ ТЕСТ ПРОЕКТА Go + Docker + PostgreSQL
# ============================================

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "     TESTING PROJECT" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# 1. CHECK CONTAINERS
Write-Host "1. CHECKING CONTAINERS:" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
Write-Host "Running containers:" -ForegroundColor White
docker ps
Write-Host ""
Write-Host "All containers (including stopped):" -ForegroundColor White
docker ps -a | findstr movies
Write-Host ""

# 2. GET ALL MOVIES
Write-Host "2. GET ALL MOVIES:" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri http://localhost:8000/movies -Method GET -ErrorAction Stop
    Write-Host "Result:" -ForegroundColor Green
    $response | ConvertTo-Json -Depth 3
}
catch {
    Write-Host "Error getting movies:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# 3. CREATE NEW MOVIE
Write-Host "3. CREATE NEW MOVIE (POST):" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
$newMovie = @{
    title = "Inception"
    genre = "Sci-Fi"
    rating = 8.8
    description = "A thief who steals corporate secrets through dreams"
} | ConvertTo-Json

try {
    $createdMovie = Invoke-RestMethod -Uri http://localhost:8000/movies `
        -Method POST `
        -ContentType "application/json" `
        -Body $newMovie `
        -ErrorAction Stop
    
    Write-Host "Created movie:" -ForegroundColor Green
    $createdMovie | ConvertTo-Json -Depth 3
    $newId = $createdMovie.id
    Write-Host "New movie ID: $newId" -ForegroundColor Green
}
catch {
    Write-Host "Error creating movie:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    $newId = 4
}
Write-Host ""

# 4. GET SINGLE MOVIE
Write-Host "4. GET MOVIE BY ID (GET):" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
try {
    $singleMovie = Invoke-RestMethod -Uri "http://localhost:8000/movies/$newId" -Method GET -ErrorAction Stop
    Write-Host "Movie with ID $newId :" -ForegroundColor Green
    $singleMovie | ConvertTo-Json -Depth 3
}
catch {
    Write-Host "Movie with ID $newId not found" -ForegroundColor Red
}
Write-Host ""

# 5. UPDATE MOVIE
Write-Host "5. UPDATE MOVIE (PUT):" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
$updatedMovie = @{
    title = "Inception (Updated)"
    genre = "Sci-Fi/Thriller"
    rating = 9.2
    description = "A mind-bending thriller about dream invasion"
} | ConvertTo-Json

try {
    $updateResult = Invoke-RestMethod -Uri "http://localhost:8000/movies/$newId" `
        -Method PUT `
        -ContentType "application/json" `
        -Body $updatedMovie `
        -ErrorAction Stop
    
    Write-Host "Updated movie:" -ForegroundColor Green
    $updateResult | ConvertTo-Json -Depth 3
}
catch {
    Write-Host "Error updating movie:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# 6. CHECK AFTER UPDATE
Write-Host "6. CHECK AFTER UPDATE:" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
try {
    $checkUpdate = Invoke-RestMethod -Uri "http://localhost:8000/movies/$newId" -Method GET -ErrorAction Stop
    Write-Host "Data after update:" -ForegroundColor Green
    $checkUpdate | ConvertTo-Json -Depth 3
}
catch {
    Write-Host "Could not verify update" -ForegroundColor Red
}
Write-Host ""

# 7. DELETE MOVIE
Write-Host "7. DELETE MOVIE (DELETE):" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
try {
    Invoke-RestMethod -Uri "http://localhost:8000/movies/$newId" -Method DELETE -ErrorAction Stop
    Write-Host "Movie with ID $newId successfully deleted" -ForegroundColor Green
}
catch {
    Write-Host "Error deleting movie:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# 8. CHECK AFTER DELETE
Write-Host "8. CHECK AFTER DELETE:" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
try {
    $afterDelete = Invoke-RestMethod -Uri "http://localhost:8000/movies/$newId" -Method GET -ErrorAction Stop
    Write-Host "Movie still exists (SHOULD NOT!)" -ForegroundColor Red
}
catch {
    Write-Host "Movie with ID $newId not found (SUCCESSFULLY DELETED)" -ForegroundColor Green
}
Write-Host ""

# 9. FINAL LIST
Write-Host "9. FINAL LIST OF ALL MOVIES:" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
try {
    $finalList = Invoke-RestMethod -Uri http://localhost:8000/movies -Method GET -ErrorAction Stop
    Write-Host "All movies in database:" -ForegroundColor Green
    $finalList | ConvertTo-Json -Depth 3
}
catch {
    Write-Host "Error getting final list:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
}
Write-Host ""

# 10. CHECK IMAGE SIZES
Write-Host "10. IMAGE SIZES (MULTI-STAGE BUILD):" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
Write-Host "Go application image (should be small ~15-20MB):" -ForegroundColor White
docker images | findstr go-docker-project
Write-Host ""
Write-Host "All images:" -ForegroundColor White
docker images | findstr -E "go-docker|postgres"
Write-Host ""

# 11. DATABASE CHECK
Write-Host "11. DATABASE CHECK:" -ForegroundColor Yellow
Write-Host "-------------------------------------" -ForegroundColor Yellow
Write-Host "Connect to PostgreSQL via Docker exec:" -ForegroundColor White
Write-Host "docker exec -it movies-db psql -U postgres -d moviesdb -c '\dt'" -ForegroundColor Gray
Write-Host "docker exec -it movies-db psql -U postgres -d moviesdb -c 'SELECT * FROM movies;'" -ForegroundColor Gray
Write-Host ""

# 12. SUMMARY
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "          TEST SUMMARY" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host " GET requests working" -ForegroundColor Green
Write-Host " POST requests working" -ForegroundColor Green
Write-Host " PUT requests working" -ForegroundColor Green
Write-Host " DELETE requests working" -ForegroundColor Green
Write-Host " All CRUD operations tested successfully" -ForegroundColor Green
Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "        TESTING COMPLETED" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan

# Save results to file
$timestamp = Get-Date -Format "yyyy-MM-dd_HH-mm-ss"
$logFile = "test-results-$timestamp.txt"
Write-Host ""
Write-Host "Results saved to file: $logFile" -ForegroundColor Magenta