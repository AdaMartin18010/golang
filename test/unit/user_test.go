package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	appuser "github.com/yourusername/golang/internal/app/user"
	domainuser "github.com/yourusername/golang/internal/domain/user"
)

// MockRepository 模拟仓储
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*domainuser.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainuser.User), args.Error(1)
}

func (m *MockRepository) Save(ctx context.Context, user *domainuser.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Update(ctx context.Context, user *domainuser.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]*domainuser.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domainuser.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		userName string
		setup    func(*MockRepository)
		wantErr  bool
	}{
		{
			name:     "success",
			email:    "test@example.com",
			userName: "Test User",
			setup: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, domainuser.ErrUserNotFound)
				m.On("Save", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "user already exists",
			email:    "existing@example.com",
			userName: "Existing User",
			setup: func(m *MockRepository) {
				existingUser := domainuser.NewUser("existing@example.com", "Existing User")
				m.On("FindByEmail", mock.Anything, "existing@example.com").Return(existingUser, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.setup(mockRepo)

			service := appuser.NewService(mockRepo)
			ctx := context.Background()

			user, err := service.CreateUser(ctx, tt.email, tt.userName)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
				assert.Equal(t, tt.userName, user.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockRepository)
	expectedUser := domainuser.NewUser("test@example.com", "Test User")
	expectedUser.ID = "1"
	expectedUser.CreatedAt = time.Now()
	expectedUser.UpdatedAt = time.Now()

	mockRepo.On("FindByID", mock.Anything, "1").Return(expectedUser, nil)

	service := appuser.NewService(mockRepo)
	ctx := context.Background()

	user, err := service.GetUser(ctx, "1")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.Name, user.Name)

	mockRepo.AssertExpectations(t)
}

