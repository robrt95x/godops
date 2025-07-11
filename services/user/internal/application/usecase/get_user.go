package usecase

import (
	"github.com/robrt95x/godops/services/user/internal/domain/entity"
	"github.com/robrt95x/godops/services/user/internal/domain/service"
)

type GetUserUseCase struct {
	userService *service.UserService
}

func NewGetUserUseCase(userService *service.UserService) *GetUserUseCase {
	return &GetUserUseCase{
		userService: userService,
	}
}

func (uc *GetUserUseCase) Execute(id string) (*entity.User, error) {
	return uc.userService.GetUserByID(id)
}

func (uc *GetUserUseCase) ExecuteGetAll() ([]*entity.User, error) {
	return uc.userService.GetAllUsers()
}
