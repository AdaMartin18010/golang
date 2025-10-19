package implementations

import (
	"errors"
	"sync"
	"time"

	"domain"
)

// MemoryUserRepository 内存实现的用户仓储
type MemoryUserRepository struct {
	users map[string]*domain.User
	mutex sync.RWMutex
}

// NewMemoryUserRepository 创建新的内存用户仓储
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// FindByID 根据ID查找用户
func (r *MemoryUserRepository) FindByID(id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *MemoryUserRepository) FindByEmail(email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

// FindAll 查找所有用户
func (r *MemoryUserRepository) FindAll() ([]*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	return users, nil
}

// FindByAgeRange 根据年龄范围查找用户
func (r *MemoryUserRepository) FindByAgeRange(minAge, maxAge int) ([]*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var users []*domain.User
	for _, user := range r.users {
		if user.Age >= minAge && user.Age <= maxAge {
			users = append(users, user)
		}
	}

	return users, nil
}

// Save 保存用户
func (r *MemoryUserRepository) Save(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查邮箱是否已存在
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return domain.ErrUserAlreadyExists
		}
	}

	// 生成ID（在实际项目中应该使用UUID）
	if user.ID == "" {
		user.ID = generateID()
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	r.users[user.ID] = user
	return nil
}

// Update 更新用户
func (r *MemoryUserRepository) Update(user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if user.ID == "" {
		return errors.New("user ID cannot be empty")
	}

	existingUser, exists := r.users[user.ID]
	if !exists {
		return domain.ErrUserNotFound
	}

	// 检查邮箱是否被其他用户使用
	for id, otherUser := range r.users {
		if id != user.ID && otherUser.Email == user.Email {
			return domain.ErrUserAlreadyExists
		}
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return nil
}

// Delete 删除用户
func (r *MemoryUserRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

// generateID 生成简单的ID（仅用于演示）
func generateID() string {
	return time.Now().Format("20060102150405")
}
