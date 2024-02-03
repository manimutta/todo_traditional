package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Task represents a task in the TODO app
type Task struct {
	gorm.Model
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	Notes     string `json:"notes"`
}

var db *gorm.DB

func main() {
	db, err := gorm.Open("sqlite3", "tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Auto Migrate the schema
	dbErr := db.AutoMigrate(&Task{})
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	router := mux.NewRouter()

	router.HandleFunc("/tasks", GetTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", GetTask).Methods("GET")
	router.HandleFunc("/tasks", CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", router)
}

// GetTasks returns all tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	var taskList []Task
	db.Find(&taskList)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(taskList)
}

// GetTask returns a single task by ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task Task
	result := db.First(&task, params["id"])

	if result.Error == gorm.ErrRecordNotFound {
		http.NotFound(w, r)
		return
	} else if result.Error != nil {
		log.Fatal(result.Error)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// CreateTask adds a new task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	json.NewDecoder(r.Body).Decode(&newTask)

	db.Create(&newTask)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newTask)
}

// UpdateTask updates a task by ID
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task Task
	result := db.First(&task, params["id"])

	if result.Error == gorm.ErrRecordNotFound {
		http.NotFound(w, r)
		return
	} else if result.Error != nil {
		log.Fatal(result.Error)
	}

	var updatedTask Task
	json.NewDecoder(r.Body).Decode(&updatedTask)

	// Update task fields
	task.Title = updatedTask.Title
	task.Completed = updatedTask.Completed
	task.Notes = updatedTask.Notes

	db.Save(&task)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DeleteTask deletes a task by ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var task Task
	result := db.First(&task, params["id"])

	if result.Error == gorm.ErrRecordNotFound {
		http.NotFound(w, r)
		return
	} else if result.Error != nil {
		log.Fatal(result.Error)
	}

	db.Delete(&task)

	w.WriteHeader(http.StatusNoContent)
}
