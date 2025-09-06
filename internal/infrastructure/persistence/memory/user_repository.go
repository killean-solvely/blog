package memory

import (
	"errors"
	"sync"

	"blog/internal/domain"
)

type UserRepository struct {
	mu    sync.RWMutex
	users map[domain.UserID]domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: map[domain.UserID]domain.User{},
	}
}

func (r *UserRepository) All() ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := []domain.User{}
	for k := range r.users {
		users = append(users, r.users[k])
	}

	return users, nil
}

func (r *UserRepository) FindByID(id domain.UserID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("no rows")
	}

	return &user, nil
}

func (r *UserRepository) Exists(id domain.UserID) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.users[id]
	return exists, nil
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, v := range r.users {
		if v.Username() == username {
			return true, nil
		}
	}

	return false, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, v := range r.users {
		if v.Email() == email {
			return true, nil
		}
	}

	return false, nil
}

func (r *UserRepository) Create(user *domain.User) (*domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.GetID()] = *user

	u := r.users[user.GetID()]
	return &u, nil
}

func (r *UserRepository) UpdateRoles(id domain.UserID, roles []domain.UserRole) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u := r.users[id]
	for _, r := range roles {
		u.AddRole(r)
	}
	r.users[id] = u

	return nil
}

func (r *UserRepository) UpdateDescription(id domain.UserID, newDescription string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u := r.users[id]
	u.UpdateDescription(newDescription)
	r.users[id] = u

	return nil
}

func (r *UserRepository) UpdatePasswordHash(id domain.UserID, newPasswordHash string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u := r.users[id]
	u.UpdatePasswordHash(newPasswordHash)
	r.users[id] = u

	return nil
}
