package usecase

import (
	"practice-7/internal/entity"
	"practice-7/internal/usecase/repo"
	"practice-7/utils"
	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (uc *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	uc.repo.Create(user)
	return user, uuid.New().String(), nil
}

func (uc *UserUseCase) LoginUser(dto *entity.LoginUserDTO) (string, error) {
	user, _ := uc.repo.FindByUsername(dto.Username)
	utils.CheckPassword(user.Password, dto.Password)
	return utils.GenerateJWT(user.ID, user.Role)
}

func (uc *UserUseCase) GetUserByID(id string) (*entity.User, error) {
	return uc.repo.FindByID(id)
}

func (uc *UserUseCase) PromoteUser(id string) error {
	return uc.repo.UpdateRole(id, "admin")
}