package usecase

import (
    "fmt"
    "my-golang-project/internal/repository"
    "my-golang-project/pkg/modules"
)

type UserUsecase struct {
    repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
    return &UserUsecase{repo: repo}
}

func (u *UserUsecase) GetUsers() ([]modules.User, error) {
    return u.repo.GetUsers()
}

func (u *UserUsecase) GetUserByID(id int) (*modules.User, error) {
    return u.repo.GetUserByID(id)
}

//Новый метод: получить только активного пользователя
func (u *UserUsecase) GetActiveUserByID(id int) (*modules.User, error) {
    return u.repo.GetActiveUserByID(id)
}

func (u *UserUsecase) CreateUser(name, email string, age *int) (int, error) {
    if name == "" {
        return 0, fmt.Errorf("имя не может быть пустым")
    }
    if email == "" {
        return 0, fmt.Errorf("email не может быть пустым")
    }
    return u.repo.CreateUser(name, email, age)
}

func (u *UserUsecase) UpdateUser(id int, name, email string, age *int) error {
    return u.repo.UpdateUser(id, name, email, age)
}

func (u *UserUsecase) DeleteUser(id int) error {
    return u.repo.DeleteUser(id)
}

//Новые методы для работы с удаленными
func (u *UserUsecase) HardDeleteUser(id int) error {
    return u.repo.HardDeleteUser(id)
}

func (u *UserUsecase) GetDeletedUsers() ([]modules.User, error) {
    return u.repo.GetDeletedUsers()
}

func (u *UserUsecase) RestoreUser(id int) error {
    return u.repo.RestoreUser(id)
}