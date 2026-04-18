package service

import (
	"errors"
	"practice8/repository"
	"testing"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("User already exists", func(t *testing.T) {
		user := &repository.User{ID: 2, Name: "Test"}
		mockRepo.EXPECT().GetByEmail("existing@test.com").Return(user, nil)

		err := userService.RegisterUser(user, "existing@test.com")
		assert.ErrorContains(t, err, "user s takim email uzhe est")
	})

	t.Run("New User -> Success", func(t *testing.T) {
		user := &repository.User{ID: 3, Name: "New"}
		mockRepo.EXPECT().GetByEmail("new@test.com").Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(nil)

		err := userService.RegisterUser(user, "new@test.com")
		assert.NoError(t, err)
	})

	t.Run("Repository error on GetByEmail", func(t *testing.T) {
		user := &repository.User{ID: 4, Name: "Error"}
		mockRepo.EXPECT().GetByEmail("error@test.com").Return(nil, errors.New("db error"))

		err := userService.RegisterUser(user, "error@test.com")
		assert.ErrorContains(t, err, "error pri proverke emaila")
	})

	t.Run("Repository error on CreateUser", func(t *testing.T) {
		user := &repository.User{ID: 5, Name: "CreateFail"}
		mockRepo.EXPECT().GetByEmail("fail@test.com").Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(errors.New("create failed"))

		err := userService.RegisterUser(user, "fail@test.com")
		assert.ErrorContains(t, err, "create failed")
	})
}

func TestUpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("Empty name", func(t *testing.T) {
		err := userService.UpdateUserName(1, "")
		assert.ErrorContains(t, err, "imya ne mozhet bit pustym")
	})

	t.Run("User not found", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(99).Return(nil, errors.New("not found"))

		err := userService.UpdateUserName(99, "NewName")
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("Successful update", func(t *testing.T) {
		user := &repository.User{ID: 10, Name: "OldName"}
		mockRepo.EXPECT().GetUserByID(10).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(user).Return(nil)

		err := userService.UpdateUserName(10, "NewName")
		assert.NoError(t, err)
		assert.Equal(t, "NewName", user.Name)
	})

	t.Run("UpdateUser fails", func(t *testing.T) {
		user := &repository.User{ID: 11, Name: "Old"}
		mockRepo.EXPECT().GetUserByID(11).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(user).Return(errors.New("update failed"))

		err := userService.UpdateUserName(11, "New")
		assert.ErrorContains(t, err, "update failed")
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("Attempt to delete admin", func(t *testing.T) {
		err := userService.DeleteUser(1)
		assert.ErrorContains(t, err, "nelzya udalit admina bratan")
	})

	t.Run("Successful delete", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(42).Return(nil)

		err := userService.DeleteUser(42)
		assert.NoError(t, err)
	})

	t.Run("Repository error", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(100).Return(errors.New("db error"))

		err := userService.DeleteUser(100)
		assert.ErrorContains(t, err, "db error")
	})
}