package port

import "github.com/robrt95x/godops/services/user/internal/domain/entity"

type UserRepository interface {
	Save(user *entity.User) error
	GetByID(id string) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
}
