package unit

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	appuser "github.com/yourusername/golang/internal/application/user"
	"github.com/yourusername/golang/internal/domain/user"
)

// MockRepository 模拟仓储
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) List(ctx context.Context, limit, offset int) ([]*user.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		req     appuser.CreateUserRequest
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name: "success",
			req: appuser.CreateUserRequest{
				Email: "test@example.com",
				Name:  "Test User",
			},
			setup: func(m *MockRepository) {
				m.On("FindByEmail", mock.Anything, "test@example.com").Return(nil, user.ErrUserNotFound)
				m.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "user already exists",
			req: appuser.CreateUserRequest{
				Email: "existing@example.com",
				Name:  "Existing User",
			},
			setup: func(m *MockRepository) {
				existingUser := &user.User{
					ID:    "1",
					Email: "existing@example.com",
					Name:  "Existing User",
				}
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

			resp, err := service.CreateUser(ctx, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.req.Email, resp.Email)
				assert.Equal(t, tt.req.Name, resp.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockRepository)
	expectedUser := &user.User{
		ID:        "1",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, "1").Return(expectedUser, nil)

	service := appuser.NewService(mockRepo)
	ctx := context.Background()

	resp, err := service.GetUser(ctx, "1")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedUser.ID, resp.ID)
	assert.Equal(t, expectedUser.Email, resp.Email)
	assert.Equal(t, expectedUser.Name, resp.Name)

	mockRepo.AssertExpectations(t)
}
