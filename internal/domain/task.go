package domain

import (
	"time"
	"workmate/internal/entity"
)

type Task struct {
	entity.Task
}

func (t *Task) StartProcessing() {
	t.Status = entity.StatusRunning
	t.UpdatedAt = time.Now().UTC()
}

func (t *Task) Complete(result string) {
	t.Status = entity.StatusCompleted
	t.Result = result
	t.UpdatedAt = time.Now().UTC()
}

func (t *Task) Fail(err error) {
	t.Status = entity.StatusFailed
	t.Error = err.Error()
	t.UpdatedAt = time.Now().UTC()
}
