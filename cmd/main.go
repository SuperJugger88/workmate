package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"workmate/internal/infrastructure"

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
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close DB:", err)
		}
	}(db)

	// Регистрируем метрики
	prometheus.MustRegister(pkg.TaskStatus)
	prometheus.MustRegister(pkg.RequestDuration)
	prometheus.MustRegister(pkg.HealthStatus)

	// Инициализация слоев
	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	handler := api.NewTaskHandler(taskService)

	router := infrastructure.SetupRouter(handler)

	// Сервер
	server := &http.Server{
		Addr:    ":" + os.Getenv("APP_PORT"),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
