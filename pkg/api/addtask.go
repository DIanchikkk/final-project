package api

import (
	"encoding/json"
	"go1f/pkg/db"
	"net/http"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка чтения JSON"})
		return
	}

	if task.Title == "" {
		WriteJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Неверный формат даты"})
		return
	}

	todayStr := now.Format("20060102")
	today, _ := time.Parse("20060102", todayStr)

	if t.Before(today) {
		if task.Repeat != "" {
			next, err := NextDate(today, task.Date, task.Repeat)
			if err != nil {
				WriteJson(w, map[string]string{"error": err.Error()})
				return
			}
			task.Date = next
		} else {
			task.Date = todayStr
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		WriteJson(w, map[string]string{"error": "Ошибка добавления задачи в базу"})
		return
	}

	WriteJson(w, map[string]interface{}{"id": id})
}
