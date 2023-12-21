package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kamalkish0r/GO-REST/pkg/db"
)

// extract id from url parameters
func GetID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, err
	}

	return id, err
}

// get all tasks
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Initialise connection to db
	database, err := db.InitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// get all the tasks from database
	tasks, err := db.GetAllTasks(database)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert tasks to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// get task by id
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := GetID(r)
	if err != nil {
		http.Error(w, "Invalid task id", http.StatusBadRequest)
		return
	}

	// initialise db connection
	database, err := db.InitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer database.Close()

	task, err := db.GetTaskByID(database, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert the task to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// add a task
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask db.Task

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// initialise db connection
	database, err := db.InitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer database.Close()

	newTaskId, err := db.CreateTask(database, newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// provide task created response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Task added successfully with ID : " + strconv.Itoa(newTaskId)))
}

// delete a task by id
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := GetID(r)
	if err != nil {
		http.Error(w, "Invalid task id", http.StatusBadRequest)
		return
	}

	// initialise connection to db
	database, err := db.InitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// delete the task
	err = db.DeleteTaskByID(database, id)
	if err != nil {
		if err == db.ErrTaskNotFound {
			http.Error(w, db.ErrTaskNotFound.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// provide task deleted response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task with ID : " + strconv.Itoa(id) + " deleted successfully."))
}

// update a task by modifying its status
func UpdateTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := GetID(r)
	if err != nil {
		http.Error(w, "Invalid task id", http.StatusBadRequest)
		return
	}

	// get new status
	var newStatus struct {
		Status string `json:"status"`
	}
	err = json.NewDecoder(r.Body).Decode(&newStatus)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// initialise db connection
	database, err := db.InitDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer database.Close()

	// update task
	err = db.UpdateTaskStatusByID(database, id, newStatus.Status)
	if err != nil {
		if err == db.ErrTaskNotFound {
			http.Error(w, "Invalid task ID", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedTask, err := db.GetTaskByID(database, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set the response content type to JSON and write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}
