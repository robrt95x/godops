package usecase

import (
	"github.com/robrt95x/godops/services/user/internal/domain/entity"
	"github.com/robrt95x/godops/services/user/internal/domain/service"
)

type CreateUserUseCase struct {
	userService *service.UserService
}

func NewCreateUserUseCase(userService *service.UserService) *CreateUserUseCase {
	return &CreateUserUseCase{
		userService: userService,
	}
}

func (uc *CreateUserUseCase) Execute(name, email string) (*entity.User, error) {
	return uc.userService.CreateUser(name, email)
}
