package service

import (
	"context"
	"time"

	"workmate/internal/domain"
	"workmate/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context) (*domain.Task, error) {
	task := domain.NewTask()
	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	go s.processTask(task.ID)
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) processTask(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	task, _ := s.GetTask(ctx, id)
	task.Status = domain.StatusRunning
	_ = s.repo.Update(ctx, task)

	// Имитация долгой операции
	select {
	case <-time.After(3 * time.Minute):
		task.Status = domain.StatusCompleted
		task.Result = "Operation completed"
	case <-ctx.Done():
		task.Status = domain.StatusFailed
		task.Error = ctx.Err().Error()
	}

	_ = s.repo.Update(ctx, task)
}
