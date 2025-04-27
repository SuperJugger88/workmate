package repository

import (
	"context"
	"database/sql"
	"log"
	"time"
	"workmate/internal/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entity.TaskEntity) error
	GetByID(ctx context.Context, id string) (*entity.TaskEntity, error)
	Update(ctx context.Context, task *entity.TaskEntity) error
	ListByStatus(ctx context.Context, status entity.TaskStatus) ([]*entity.TaskEntity, error)
}

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) Create(ctx context.Context, task *entity.TaskEntity) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO tasks (id, status, created_at, updated_at) VALUES ($1, $2, $3, $4)",
		task.ID, task.Status, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

func (r *PostgresTaskRepository) GetByID(ctx context.Context, id string) (*entity.TaskEntity, error) {
	task := &entity.TaskEntity{}
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

func (r *PostgresTaskRepository) Update(ctx context.Context, task *entity.TaskEntity) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET status = $1, result = $2, error = $3, updated_at = $4 WHERE id = $5`,
		task.Status,
		task.Result,
		task.Error,
		time.Now().UTC(),
		task.ID,
	)
	return err
}

func (r *PostgresTaskRepository) ListByStatus(ctx context.Context, status entity.TaskStatus) ([]*entity.TaskEntity, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, status, result, error, created_at, updated_at FROM tasks WHERE status = $1",
		status,
	)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatalf("failed to close rows: %v", err)
		}
	}(rows)

	var tasks []*entity.TaskEntity
	for rows.Next() {
		task := &entity.TaskEntity{}
		if err := rows.Scan(
			&task.ID,
			&task.Status,
			&task.Result,
			&task.Error,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
