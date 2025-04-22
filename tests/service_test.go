package tests

import (
	"context"
	"testing"
	"time"

	"workmate/internal/entity"
	"workmate/internal/repository/mocks"
	"workmate/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskService_GetTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	svc := service.NewTaskService(mockRepo)

	expectedTask := &entity.Task{
		ID:        "test-id",
		Status:    entity.StatusCompleted,
		CreatedAt: time.Now(),
	}

	mockRepo.On("GetByID", mock.Anything, "test-id").
		Return(expectedTask, nil)

	task, err := svc.GetTask(context.Background(), "test-id")
	assert.NoError(t, err)
	assert.Equal(t, expectedTask.ID, task.ID)
	mockRepo.AssertExpectations(t)
}
