package db

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

var ErrTaskNotFound = errors.New("Task not found")

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS Task (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		description TEXT,
		status VARCHAR(20) DEFAULT 'Pending' CHECK (status IN ('Pending', 'InProgress', 'Complete'))
	)
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func CreateTask(db *sql.DB, newTask Task) (int, error) {
	var assignedID int
	err := db.QueryRow(
		"INSERT INTO Task (title, description, status) VALUES ($1, $2, $3) RETURNING id",
		newTask.Title, newTask.Description, newTask.Status,
	).Scan(&assignedID)
	if err != nil {
		return -1, err
	}

	return assignedID, nil
}

func GetAllTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT * FROM Task")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTaskByID(db *sql.DB, id int) (Task, error) {
	row := db.QueryRow("SELECT * FROM Task WHERE Task.id = $1", id)

	var task Task
	err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskNotFound
		}
		return Task{}, err
	}
	return task, nil
}

func DeleteTaskByID(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM Task WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func UpdateTaskStatusByID(db *sql.DB, id int, status string) error {
	result, err := db.Exec("UPDATE Task SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	return nil
}
