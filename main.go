package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var tasks = []Task{}
var tasksMu sync.Mutex
var nextID = 1

var ready atomic.Bool

var requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "requests_total",
	Help: "Total number of requests handled.",
})

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	requestsTotal.Inc()
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	requestsTotal.Inc()
	tasksMu.Lock()
	snapshot := make([]Task, len(tasks))
	copy(snapshot, tasks)
	tasksMu.Unlock()

	writeJSON(w, http.StatusOK, snapshot)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	requestsTotal.Inc()
	var body struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if body.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	tasksMu.Lock()
	task := Task{ID: nextID, Title: body.Title}
	nextID++
	tasks = append(tasks, task)
	tasksMu.Unlock()

	writeJSON(w, http.StatusCreated, task)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	if !ready.Load() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not ready"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func main() {
	go func() {
		time.Sleep(10 * time.Second)
		ready.Store(true)
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /tasks", getTasksHandler)
	mux.HandleFunc("POST /tasks", createTaskHandler)
	mux.HandleFunc("GET /ready", readyHandler)
	mux.Handle("GET /metrics", promhttp.Handler())

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
