package service

import (
	"ta/internal/domain"
	"ta/internal/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

var _ repository.UserRepositoryInterface = (*MockUserRepository)(nil)

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		name    string
		user    *domain.User
		setup   func(*MockUserRepository)
		wantErr bool
		errMsg  string
	}{
		{
			name: "successful registration",
			user: &domain.User{
				ID:        uuid.New(),
				Firstname: "John",
				Lastname:  "Doe",
				Age:       25,
				Password:  "password123",
			},
			setup: func(m *MockUserRepository) {
				m.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "age less than 18",
			user: &domain.User{
				Age:      17,
				Password: "password123",
			},
			setup:   func(m *MockUserRepository) {},
			wantErr: true,
			errMsg:  "user must be at least 18 years old",
		},
		{
			name: "password too short",
			user: &domain.User{
				Age:      25,
				Password: "short",
			},
			setup:   func(m *MockUserRepository) {},
			wantErr: true,
			errMsg:  "password must be at least 8 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			tt.setup(mockRepo)
			service := NewUserService(mockRepo)

			err := service.Register(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

