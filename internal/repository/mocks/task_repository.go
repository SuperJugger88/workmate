package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"workmate/internal/entity"
)

// MockTaskRepository реализует TaskRepository для тестирования
type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *entity.TaskEntity) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*entity.TaskEntity, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.TaskEntity), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *entity.TaskEntity) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) ListByStatus(ctx context.Context, status entity.TaskStatus) ([]*entity.TaskEntity, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]*entity.TaskEntity), args.Error(1)
}
