package repository_test

import (
	"context"
	"testing"

	"workmate/internal/domain"
	"workmate/internal/repository"
	"workmate/pkg"

	"github.com/stretchr/testify/assert"
)

func TestTaskRepository(t *testing.T) {
	// Инициализация тестовой БД
	db, err := pkg.InitDB()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTaskRepository(db)
	ctx := context.Background()

	// Тест создания задачи
	t.Run("Create and Get Task", func(t *testing.T) {
		task := domain.NewTask()
		err := repo.Create(ctx, task)
		assert.NoError(t, err)

		found, err := repo.GetByID(ctx, task.ID)
		assert.NoError(t, err)
		assert.Equal(t, task.ID, found.ID)
	})
}
