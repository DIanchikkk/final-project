package api

import (
	"encoding/json"
	"net/http"
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

		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			WriteJson(w, map[string]string{"error": "Ошибка чтения JSON"})
			return
		}

		if task.ID == "" {
			WriteJson(w, map[string]string{"error": "Не указан идентификатор задачи"})
			return
		}
		if task.Title == "" {
			WriteJson(w, map[string]string{"error": "Не указан заголовок задачи"})
			return
		}
		if _, err := time.Parse("20060102", task.Date); err != nil {
			WriteJson(w, map[string]string{"error": "Неверный формат даты"})
			return
		}

		if err := db.UpdateTask(&task); err != nil {
			WriteJson(w, map[string]string{"error": err.Error()})
			return
		}

		WriteJson(w, map[string]interface{}{})

	default:
		WriteJson(w, map[string]string{"error": "Метод не поддерживается"})
	}
}
