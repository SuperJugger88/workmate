package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"workmate/internal/entity"
	"workmate/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask godoc
// @Summary Create new async task
// @Description Create new long-running task
// @Tags tasks
// @Accept  json
// @Produce  json
// @Success 201 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task, err := h.service.CreateTask(r.Context())
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"task_id": task.ID})
	if err != nil {
		log.Fatalf("Failed to encode task response: %v", err)
	}
}

// GetTaskStatus godoc
// @Summary Get task status
// @Description Get task status by ID
// @Tags tasks
// @Produce  json
// @Param id path string true "Task ID"
// @Success 200 {object} domain.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTaskStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID := vars["id"]

	task, err := h.service.GetTask(r.Context(), taskID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		log.Fatalf("Failed to encode task response: %v", err)
	}
}

// ListTasksByStatus godoc
// @Summary List tasks by status
// @Description Get tasks filtered by status
// @Tags tasks
// @Produce  json
// @Param status path string true "Task status (pending, running, completed, failed)"
// @Success 200 {array} domain.Task
// @Failure 400 {object} map[string]string
// @Router /tasks/status/{status} [get]
func (h *TaskHandler) ListTasksByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := entity.TaskStatus(vars["status"])

	if !isValidStatus(status) {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	tasks, err := h.service.ListTasksByStatus(r.Context(), status)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		log.Fatalf("Failed to encode task response: %v", err)
	}
}

func isValidStatus(s entity.TaskStatus) bool {
	switch s {
	case entity.StatusPending, entity.StatusRunning,
		entity.StatusCompleted, entity.StatusFailed:
		return true
	}
	return false
}
