package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var tasks = []Task{}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, tasks)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /tasks", getTasksHandler)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
