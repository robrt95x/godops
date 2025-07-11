package repository

import (
	"sync"

	"github.com/robrt95x/godops/services/user/internal/domain/entity"
)

type MemoryUserRepository struct {
	users map[string]*entity.User
	mutex sync.RWMutex
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (r *MemoryUserRepository) Save(user *entity.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.users[user.ID] = user
	return nil
}

func (r *MemoryUserRepository) GetByID(id string) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	
	return user, nil
}

func (r *MemoryUserRepository) GetByEmail(email string) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	
	return nil, nil
}

func (r *MemoryUserRepository) GetAll() ([]*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	
	return users, nil
}
