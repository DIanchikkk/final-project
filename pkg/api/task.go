package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go1f/pkg/db"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		AddTaskHandler(w, r)

	case http.MethodGet:
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

		WriteJson(w, task)

	case http.MethodPut:
		var task db.Task

		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			WriteJson(w, map[string]string{"error": "Ошибка чтения JSON"})
			return
		}

		if task.ID == "" || task.Title == "" {
			WriteJson(w, map[string]string{"error": "Некорректные данные задачи"})
			return
		}

		if task.Comment != "" &&
			task.Comment != strings.TrimSpace(task.Comment) {
			WriteJson(w, map[string]string{"error": "Некорректные данные задачи"})
			return
		}

		if _, err := time.Parse("20060102", task.Date); err != nil {
			WriteJson(w, map[string]string{"error": "Неверный формат даты"})
			return
		}

		// ⚠️ ОСТАВЛЯЕМ: важно для TestEditTask
		if task.Repeat != "" {
			now := time.Now()
			if _, err := NextDate(now, now.Format("20060102"), task.Repeat); err != nil {
				WriteJson(w, map[string]string{"error": err.Error()})
				return
			}
		}

		task.Title = strings.TrimSpace(task.Title)
		task.Comment = strings.TrimSpace(task.Comment)

		if err := db.UpdateTask(&task); err != nil {
			WriteJson(w, map[string]string{"error": "Задача не найдена"})
			return
		}

		WriteJson(w, map[string]any{})

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			WriteJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
			return
		}

		if err := db.DeleteTask(id); err != nil {
			WriteJson(w, map[string]string{"error": "Задача не найдена"})
			return
		}

		WriteJson(w, map[string]any{})

	default:
		WriteJson(w, map[string]string{"error": "Метод не поддерживается"})
	}
}
