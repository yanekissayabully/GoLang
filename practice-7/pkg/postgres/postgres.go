package postgres

import (
	"practice-7/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	Conn *gorm.DB
}

func New() *Postgres {
	dsn := "host=localhost user=myuser password=mypass dbname=practice7 port=5432 sslmode=disable"
	conn, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	conn.AutoMigrate(&entity.User{})
	return &Postgres{Conn: conn}
}