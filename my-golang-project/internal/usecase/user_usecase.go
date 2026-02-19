package usecase

import (
    "fmt"
    "my-golang-project/internal/repository"
    "my-golang-project/pkg/modules"
)

type UserUsecase struct {
    repo repository.UserRepository // Зависимость от интерфейса, а не от конкретной реализации!
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

func (u *UserUsecase) CreateUser(name, email string, age *int) (int, error) {
    // Здесь можно добавить бизнес-логику, например, валидацию email или возраст>0
    if name == "" {
        return 0, fmt.Errorf("имя не может быть пустым")
    }
    if email == "" {
        return 0, fmt.Errorf("email не может быть пустым")
    }
    // Валидация формата email - сложнее, пока пропустим
    return u.repo.CreateUser(name, email, age)
}

func (u *UserUsecase) UpdateUser(id int, name, email string, age *int) error {
    // Тоже можно добавить проверки
    return u.repo.UpdateUser(id, name, email, age)
}

func (u *UserUsecase) DeleteUser(id int) error {
    return u.repo.DeleteUser(id)
}