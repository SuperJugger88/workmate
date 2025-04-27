package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"workmate/internal/api"
	"workmate/internal/repository"
	"workmate/internal/service"
	"workmate/pkg"
)

// @title Task API
// @version 1.0
// @description API for managing long-running tasks
func main() {
	// Инициализация БД
	db, err := pkg.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	defer db.Close()

	// Регистрируем метрики
	prometheus.MustRegister(pkg.TaskStatus)
	prometheus.MustRegister(pkg.RequestDuration)

	// Инициализация слоев
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	handler := api.NewTaskHandler(taskService)

	// Роутер с Swagger
	router := mux.NewRouter()
	router.HandleFunc("/tasks", handler.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", handler.GetTaskStatus).Methods("GET")
	router.HandleFunc("/tasks/status/{status}", handler.ListTasksByStatus).Methods("GET")

	// Swagger
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	// Metrics
	http.Handle("/metrics", promhttp.Handler())

	// Сервер
	server := &http.Server{
		Addr:    ":" + os.Getenv("APP_PORT"),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}
}
