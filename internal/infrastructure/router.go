package infrastructure

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"workmate/internal/api"
)

func SetupRouter(handler *api.TaskHandler) *mux.Router {
	router := mux.NewRouter()

	// Api
	router.HandleFunc("/api/v1/tasks", handler.CreateTask).Methods("POST")
	router.HandleFunc("/api/v1/tasks/{id}", handler.GetTaskStatus).Methods("GET")
	router.HandleFunc("/api/v1/tasks/status/{status}", handler.ListTasksByStatus).Methods("GET")

	// Swagger
	router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	// Metrics
	router.Handle("/metrics", promhttp.Handler())

	return router
}
