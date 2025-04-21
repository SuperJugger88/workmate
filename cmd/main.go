package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"workmate/internal/api"
	"workmate/internal/repository"
	"workmate/internal/service"
	"workmate/pkg"
)

func main() {
	// Инициализация БД
	db, err := pkg.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Инициализация слоев
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	handler := api.NewTaskHandler(taskService)

	// Настройка роутера
	router := mux.NewRouter()
	router.HandleFunc("/tasks", handler.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", handler.GetTaskStatus).Methods("GET")

	// Запуск сервера
	server := &http.Server{
		Addr:    os.Getenv("APP_PORT"),
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
}
