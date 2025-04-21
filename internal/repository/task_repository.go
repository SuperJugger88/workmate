package repository

import (
	"context"
	"database/sql"
	"time"

	"workmate/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
}

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO tasks (id, status, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		task.ID, task.Status, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	task := &domain.Task{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, status, result, error, created_at, updated_at FROM tasks WHERE id = $1",
		id,
	).Scan(
		&task.ID,
		&task.Status,
		&task.Result,
		&task.Error,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	return task, err
}

func (r *PostgresTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE tasks 
		SET status = $1, result = $2, error = $3, updated_at = $4 
		WHERE id = $5`,
		task.Status,
		task.Result,
		task.Error,
		time.Now().UTC(),
		task.ID,
	)
	return err
}
