package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты: %v", err)
	}

	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	addYear := func(t time.Time) time.Time {
		return time.Date(t.Year()+1, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}

	switch {
	case repeat == "y":
		next := addYear(start)
		for !next.After(now) {
			next = addYear(next)
		}
		return next.Format(dateFormat), nil
	case strings.HasPrefix(repeat, "d"):
		parts := strings.Fields(repeat)
		if len(parts) != 2 || parts[0] != "d" {
			return "", fmt.Errorf("неподдерживаемый формат: %s", repeat)
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("не удалось преобразовать интервал в число: %v", err)
		}
		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("интервал должен быть от 1 до 400 дней, получено: %d", interval)
		}

		next := start.AddDate(0, 0, interval)
		for !next.After(now) {
			next = next.AddDate(0, 0, interval)
		}
		return next.Format(dateFormat), nil
	default:
		return "", fmt.Errorf("неподдерживаемый формат: %s", repeat)
	}
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dstart := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error
	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("не удалось разобрать now: %v", err), http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}
