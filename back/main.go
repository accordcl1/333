package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

var (
	tasks    []Task
	mu       sync.Mutex
	dataFile = "/data/tasks.json"
)

func loadTasks() {
	file, err := os.ReadFile(dataFile)
	if err == nil {
		_ = json.Unmarshal(file, &tasks)
	}
}

func saveTasks() {
	data, _ := json.MarshalIndent(tasks, "", "  ")
	_ = os.WriteFile(dataFile, data, 0644)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mu.Lock()
	defer mu.Unlock()

	if r.Method == http.MethodGet {
		if tasks == nil {
			tasks = []Task{}
		}
		_ = json.NewEncoder(w).Encode(tasks)
		return
	}

	if r.Method == http.MethodPost {
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil || t.Text == "" {
			http.Error(w, `{"error":"invalid input"}`, http.StatusBadRequest)
			return
		}
		t.ID = len(tasks) + 1
		tasks = append(tasks, t)
		saveTasks()
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(t)
		return
	}

	http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func main() {
	_ = os.MkdirAll("/data", 0755)
	loadTasks()

	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/health", healthHandler)

	log.Println("Backend running on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
