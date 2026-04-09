package repo

import (
	"practice-7/internal/entity"
	"practice-7/pkg/postgres"
)

type UserRepo struct {
	db *postgres.Postgres
}

func NewUserRepo(db *postgres.Postgres) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *entity.User) error {
	return r.db.Conn.Create(user).Error
}

func (r *UserRepo) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	r.db.Conn.Where("username = ?", username).First(&user)
	return &user, nil
}

func (r *UserRepo) FindByID(id string) (*entity.User, error) {
	var user entity.User
	r.db.Conn.Where("id = ?", id).First(&user)
	return &user, nil
}

func (r *UserRepo) UpdateRole(id, role string) error {
	return r.db.Conn.Model(&entity.User{}).Where("id = ?", id).Update("role", role).Error
}