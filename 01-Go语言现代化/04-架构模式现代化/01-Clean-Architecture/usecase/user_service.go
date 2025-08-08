package usecase

import (
	"domain"
	"repository"
)

// UserService 用户服务用例层
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建新的用户服务
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(email, name string, age int) (*domain.User, error) {
	// 创建用户实体
	user, err := domain.NewUser(email, name, age)
	if err != nil {
		return nil, err
	}

	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.FindByEmail(email)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// 保存用户
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

// GetUserByEmail 根据邮箱获取用户
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(email)
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() ([]*domain.User, error) {
	return s.userRepo.FindAll()
}

// UpdateUserProfile 更新用户资料
func (s *UserService) UpdateUserProfile(id, name string, age int) (*domain.User, error) {
	// 获取现有用户
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 更新用户资料
	if err := user.UpdateProfile(name, age); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

// GetAdultUsers 获取成年用户
func (s *UserService) GetAdultUsers() ([]*domain.User, error) {
	// 使用年龄范围查询获取成年用户
	return s.userRepo.FindByAgeRange(18, 150)
}

// GetUsersByAgeRange 根据年龄范围获取用户
func (s *UserService) GetUsersByAgeRange(minAge, maxAge int) ([]*domain.User, error) {
	if minAge < 0 || maxAge < 0 || minAge > maxAge {
		return nil, domain.NewBusinessError("INVALID_AGE_RANGE", "invalid age range")
	}

	return s.userRepo.FindByAgeRange(minAge, maxAge)
}
