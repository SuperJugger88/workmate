package service

import (
	"context"
	"time"

	"workmate/internal/domain"
	"workmate/internal/entity"
	"workmate/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context) (*domain.Task, error) {
	entityTask := entity.NewTask()
	if err := s.repo.Create(ctx, entityTask); err != nil {
		return nil, err
	}

	domainTask := &domain.Task{TaskEntity: *entityTask}
	go s.processTask(domainTask)
	return domainTask, nil
}

func (s *TaskService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	entityTask, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.Task{TaskEntity: *entityTask}, nil
}

func (s *TaskService) ListTasksByStatus(ctx context.Context, status entity.TaskStatus) ([]*domain.Task, error) {
	entityTasks, err := s.repo.ListByStatus(ctx, status)
	if err != nil {
		return nil, err
	}

	domainTasks := make([]*domain.Task, len(entityTasks))
	for i, t := range entityTasks {
		domainTasks[i] = &domain.Task{TaskEntity: *t}
	}
	return domainTasks, nil
}

func (s *TaskService) processTask(task *domain.Task) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	task.StartProcessing()
	_ = s.repo.Update(ctx, &task.TaskEntity)

	// Имитация долгой операции
	select {
	case <-time.After(3 * time.Minute):
		task.Complete("Operation completed")
	case <-ctx.Done():
		task.Fail(ctx.Err())
	}

	_ = s.repo.Update(ctx, &task.TaskEntity)
}
