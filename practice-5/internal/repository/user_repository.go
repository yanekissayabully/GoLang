package repository

import (
    "database/sql"
    "fmt"
    "strings"
    "time"
    "math/rand"
    "practice5/internal/models"
    "github.com/google/uuid"
    _ "github.com/lib/pq"
)

type Repository struct {
    db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) SeedDatabase() {
    r.db.Exec("DELETE FROM user_friends")
    r.db.Exec("DELETE FROM users")

    firstNames := []string{"Carti", "Travis", "Post", "Uzi", "Skies", "Pill", "Thug", "Kami", "Kid", "Kanye"}
    lastNames := []string{"Playboi", "Scott", "Malone", "Lil", "Lil", "Thrill", "Young", "Zilla", "Mad", "West"}
    domains := []string{"gmail.com", "yandex.ru", "mail.ru", "yahoo.com"}
    genders := []string{"male", "female", "other"}

    var users []uuid.UUID
    
    for i := 0; i < 30; i++ {
        name := firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
        email := fmt.Sprintf("user%d@%s", i, domains[rand.Intn(len(domains))])
        gender := genders[rand.Intn(len(genders))]
        birthDate := time.Now().AddDate(-(18 + rand.Intn(42)), -rand.Intn(12), -rand.Intn(28))
        
        var id uuid.UUID
        r.db.QueryRow(`
            INSERT INTO users (id, name, email, gender, birth_date) 
            VALUES (gen_random_uuid(), $1, $2, $3, $4) 
            RETURNING id`,
            name, email, gender, birthDate,
        ).Scan(&id)
        
        users = append(users, id)
    }

    for i := 0; i < 5; i++ {
        for j := i + 1; j < 5; j++ {
            r.addFriendship(users[i], users[j])
        }
    }
    
    for i := 5; i < 10; i++ {
        for j := i + 1; j < 10; j++ {
            r.addFriendship(users[i], users[j])
        }
    }
    
    commonFriends := []uuid.UUID{users[0], users[1], users[2], users[5], users[6]}
    
    for _, friend := range commonFriends {
        r.addFriendship(users[10], friend)
        r.addFriendship(users[11], friend)
        r.addFriendship(users[12], friend)
    }
    
    for i := 0; i < 50; i++ {
        u1 := rand.Intn(len(users))
        u2 := rand.Intn(len(users))
        if u1 != u2 {
            r.addFriendship(users[u1], users[u2])
        }
    }
}

func (r *Repository) addFriendship(userID, friendID uuid.UUID) {
    r.db.Exec(`
        INSERT INTO user_friends (user_id, friend_id) 
        VALUES ($1, $2), ($2, $1) 
        ON CONFLICT DO NOTHING`,
        userID, friendID,
    )
}

func (r *Repository) GetPaginatedUsers(page int, pageSize int, filters models.FilterParams, sort models.SortParams) (models.PaginatedResponse, error) {
    offset := (page - 1) * pageSize
    
    baseQuery := "FROM users WHERE 1=1"
    args := []interface{}{}
    argCount := 1
    
    if filters.ID != nil {
        baseQuery += fmt.Sprintf(" AND id = $%d", argCount)
        args = append(args, *filters.ID)
        argCount++
    }
    
    if filters.Name != nil {
        baseQuery += fmt.Sprintf(" AND name ILIKE $%d", argCount)
        args = append(args, "%"+*filters.Name+"%")
        argCount++
    }
    
    if filters.Email != nil {
        baseQuery += fmt.Sprintf(" AND email ILIKE $%d", argCount)
        args = append(args, "%"+*filters.Email+"%")
        argCount++
    }
    
    if filters.Gender != nil {
        baseQuery += fmt.Sprintf(" AND gender = $%d", argCount)
        args = append(args, *filters.Gender)
        argCount++
    }
    
    if filters.BirthDate != nil {
        baseQuery += fmt.Sprintf(" AND birth_date = $%d", argCount)
        args = append(args, *filters.BirthDate)
        argCount++
    }
    
    var totalCount int
    countQuery := "SELECT COUNT(*) " + baseQuery
    r.db.QueryRow(countQuery, args...).Scan(&totalCount)
    
    orderBy := "ORDER BY id"
    if sort.Field != "" {
        validFields := map[string]bool{
            "id": true, "name": true, "email": true, 
            "gender": true, "birth_date": true,
        }
        
        if validFields[sort.Field] {
            direction := "ASC"
            if strings.ToUpper(sort.Direction) == "DESC" {
                direction = "DESC"
            }
            orderBy = fmt.Sprintf("ORDER BY %s %s", sort.Field, direction)
        }
    }
    
    paginatedQuery := fmt.Sprintf("SELECT id, name, email, gender, birth_date %s %s LIMIT $%d OFFSET $%d", 
        baseQuery, orderBy, argCount, argCount+1)
    
    args = append(args, pageSize, offset)
    
    rows, _ := r.db.Query(paginatedQuery, args...)
    defer rows.Close()
    
    var users []models.User
    for rows.Next() {
        var u models.User
        rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate)
        users = append(users, u)
    }
    
    return models.PaginatedResponse{
        Data:       users,
        TotalCount: totalCount,
        Page:       page,
        PageSize:   pageSize,
    }, nil
}

func (r *Repository) GetCommonFriends(user1ID, user2ID uuid.UUID) ([]models.User, error) {
    query := `
        SELECT u.id, u.name, u.email, u.gender, u.birth_date
        FROM users u
        WHERE u.id IN (
            SELECT f1.friend_id 
            FROM user_friends f1
            INNER JOIN user_friends f2 ON f1.friend_id = f2.friend_id
            WHERE f1.user_id = $1 AND f2.user_id = $2
        )
        ORDER BY u.name
    `
    
    rows, _ := r.db.Query(query, user1ID, user2ID)
    defer rows.Close()
    
    var commonFriends []models.User
    for rows.Next() {
        var u models.User
        rows.Scan(&u.ID, &u.Name, &u.Email, &u.Gender, &u.BirthDate)
        commonFriends = append(commonFriends, u)
    }
    
    return commonFriends, nil
}