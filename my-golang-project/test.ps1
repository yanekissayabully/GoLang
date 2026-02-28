# ============================================
# TESTING ALL ENDPOINTS
# ============================================

Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TESTING MY-GOLANG-PROJECT API" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan
Write-Host ""

$API_KEY = "my-secret-key-123"
$BASE_URL = "http://localhost:8080"

# Function for making requests
function Test-Request {
    param($Method, $Url, $Body, $Description)
    
    Write-Host "▶ $Description" -ForegroundColor Yellow
    
    $params = @{
        Method = $Method
        Uri = "$BASE_URL$Url"
        Headers = @{
            "X-API-KEY" = $API_KEY
            "Content-Type" = "application/json"
        }
    }
    
    if ($Body) {
        $params.Body = ($Body | ConvertTo-Json)
    }
    
    try {
        $response = Invoke-RestMethod @params
        Write-Host "  ✅ Success:" -ForegroundColor Green
        $response | ConvertTo-Json | Write-Host
    } catch {
        Write-Host "  ❌ Error:" -ForegroundColor Red
        Write-Host "  $($_.Exception.Message)"
    }
    Write-Host ""
}

# ============================================
# TEST 1: Healthcheck (no key)
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 1: Healthcheck (no API key)" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Get -Uri "$BASE_URL/health"
    Write-Host "✅ Healthcheck OK:" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Healthcheck error:" -ForegroundColor Red
    Write-Host $_.Exception.Message
}
Write-Host ""

# ============================================
# TEST 2: Access without API key (should be 401)
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 2: Access without API key (should return 401)" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Get -Uri "$BASE_URL/users" -ErrorAction Stop
} catch {
    if ($_.Exception.Response.StatusCode -eq 401) {
        Write-Host "✅ Success: Got 401 Unauthorized" -ForegroundColor Green
    } else {
        Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    }
}
Write-Host ""

# ============================================
# TEST 3: Create users
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 3: Creating users" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

# Create first user
Test-Request -Method Post -Url "/users" -Body @{
    name = "Ivan Petrov"
    email = "ivan@mail.com"
    age = 25
} -Description "Creating user 1"

# Create second user
Test-Request -Method Post -Url "/users" -Body @{
    name = "Maria Sidorova"
    email = "maria@mail.com"
    age = 30
} -Description "Creating user 2"

# Create third user
Test-Request -Method Post -Url "/users" -Body @{
    name = "Petr Ivanov"
    email = "petr@mail.com"
    age = 35
} -Description "Creating user 3"

# ============================================
# TEST 4: Get all users
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 4: Getting all users" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Get -Uri "$BASE_URL/users" -Headers @{"X-API-KEY"=$API_KEY}
    Write-Host "✅ Total users: $($response.Count)" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 5: Get user by ID
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 5: Getting user ID=1" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Get -Uri "$BASE_URL/users/1" -Headers @{"X-API-KEY"=$API_KEY}
    Write-Host "✅ User found:" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 6: Update user
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 6: Updating user ID=2" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Put -Uri "$BASE_URL/users/2" `
        -Headers @{"X-API-KEY"=$API_KEY; "Content-Type"="application/json"} `
        -Body '{"name":"Maria Petrova","email":"masha@mail.com","age":31}'
    Write-Host "✅ Update:" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 7: Soft Delete
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 7: Soft deleting user ID=2" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Delete -Uri "$BASE_URL/users/2" -Headers @{"X-API-KEY"=$API_KEY}
    Write-Host "✅ Soft delete:" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 8: Check after delete
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 8: Active users after soft delete" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Get -Uri "$BASE_URL/users" -Headers @{"X-API-KEY"=$API_KEY}
    Write-Host "✅ Active users:" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 9: Get deleted users
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 9: Getting deleted users" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Get -Uri "$BASE_URL/users/deleted" -Headers @{"X-API-KEY"=$API_KEY}
    Write-Host "✅ Deleted users:" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 10: Restore user
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 10: Restoring user ID=2" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Post -Uri "$BASE_URL/users/2/restore" -Headers @{"X-API-KEY"=$API_KEY}
    Write-Host "✅ Restore:" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 11: Check after restore
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 11: Active users after restore" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Get -Uri "$BASE_URL/users" -Headers @{"X-API-KEY"=$API_KEY}
    Write-Host "✅ Total users: $($response.Count)" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 12: Hard Delete
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 12: Hard deleting user ID=3" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Delete -Uri "$BASE_URL/users/3/hard" -Headers @{"X-API-KEY"=$API_KEY}
    Write-Host "✅ Hard delete:" -ForegroundColor Green
    $response | ConvertTo-Json | Write-Host
} catch {
    Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# ============================================
# TEST 13: Non-existent ID
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 13: Trying to get non-existent ID=999" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Get -Uri "$BASE_URL/users/999" -Headers @{"X-API-KEY"=$API_KEY}
} catch {
    if ($_.Exception.Response.StatusCode -eq 404) {
        Write-Host "✅ Success: Got 404 Not Found" -ForegroundColor Green
    } else {
        Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
    }
}
Write-Host ""

# ============================================
# TEST 14: Duplicate email
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "TEST 14: Trying to create user with duplicate email" -ForegroundColor Cyan
Write-Host "============================================" -ForegroundColor Cyan

try {
    $response = Invoke-RestMethod -Method Post -Uri "$BASE_URL/users" `
        -Headers @{"X-API-KEY"=$API_KEY; "Content-Type"="application/json"} `
        -Body '{"name":"Duplicate","email":"ivan@mail.com","age":40}'
} catch {
    Write-Host "✅ Success: Got error about duplicate email" -ForegroundColor Green
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Yellow
}
Write-Host ""

# ============================================
# SUMMARY
# ============================================
Write-Host "============================================" -ForegroundColor Cyan
Write-Host "✅ TESTING COMPLETED!" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Cyan



