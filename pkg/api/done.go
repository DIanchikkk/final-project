package api

import (
	"net/http"
	"time"

	"go1f/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		WriteJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	// одноразовая задача
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			WriteJson(w, map[string]string{"error": "Задача не найдена"})
			return
		}
		WriteJson(w, map[string]any{})
		return
	}

	current, err := time.Parse("20060102", task.Date)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Неверный формат даты"})
		return
	}

	next, err := NextDate(current, task.Date, task.Repeat)
	if err != nil {
		WriteJson(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		WriteJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	WriteJson(w, map[string]any{})
}
