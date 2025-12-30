package api

import (
	"encoding/json"
	"net/http"

	"go1f/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		WriteJson(w, map[string]string{"error": err.Error()})
		return
	}

	WriteJson(w, TasksResp{Tasks: tasks})
}

func WriteJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
