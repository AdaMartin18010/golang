package repository

import "clean-architecture-example/domain"

// UserRepository 定义用户数据访问接口
type UserRepository interface {
	UserReader
	UserWriter
}

// UserReader 定义用户读取操作接口
type UserReader interface {
	FindByID(id string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll() ([]*domain.User, error)
	FindByAgeRange(minAge, maxAge int) ([]*domain.User, error)
}

// UserWriter 定义用户写入操作接口
type UserWriter interface {
	Save(user *domain.User) error
	Update(user *domain.User) error
	Delete(id string) error
}

// TransactionManager 定义事务管理接口
type TransactionManager interface {
	Begin() (Transaction, error)
}

// Transaction 定义事务接口
type Transaction interface {
	Commit() error
	Rollback() error
	UserRepository() UserRepository
}
