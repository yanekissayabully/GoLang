package usecase

import "practice-7/internal/entity"

type UserInterface interface {
	RegisterUser(*entity.User) (*entity.User, string, error)
	LoginUser(*entity.LoginUserDTO) (string, error)
	GetUserByID(string) (*entity.User, error)
	PromoteUser(string) error
}