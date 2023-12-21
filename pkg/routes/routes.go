package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kamalkish0r/GO-REST/pkg/controller"
)

func Routes(router *mux.Router) {
	router.HandleFunc("/tasks", controller.GetTasksHandler).Methods(http.MethodGet)
	router.HandleFunc("/tasks/{id}", controller.GetTaskHandler).Methods(http.MethodGet)
	router.HandleFunc("/tasks", controller.AddTaskHandler).Methods(http.MethodPost)
	router.HandleFunc("/tasks/{id}", controller.UpdateTaskStatusHandler).Methods(http.MethodPut)
	router.HandleFunc("/tasks/{id}", controller.DeleteTaskHandler).Methods(http.MethodDelete)
}
