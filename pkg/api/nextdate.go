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

	switch {
	case repeat == "y":
		next := start

		if next.After(now) {
			next = next.AddDate(1, 0, 0)
		}

		for !next.After(now) {
			year := next.Year() + 1
			month := start.Month()
			day := start.Day()

			if month == time.February && day == 29 {
				if !(year%4 == 0 && (year%100 != 0 || year%400 == 0)) {
					month = time.March
					day = 1
				}
			}

			next = time.Date(year, month, day, 0, 0, 0, 0, start.Location())
		}

		return next.Format(dateFormat), nil

	case strings.HasPrefix(repeat, "d"):
		parts := strings.Fields(repeat)
		if len(parts) < 2 {
			return "", fmt.Errorf("неподдерживаемый формат: %s", repeat)
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("не удалось преобразовать интервал в число: %v", err)
		}
		if interval < 1 || interval > 400 {
			return "", fmt.Errorf("интервал должен быть от 1 до 400 дней, получено: %d", interval)
		}

		next := start

		if next.After(now) {
			next = next.AddDate(0, 0, interval)
		}

		for !next.After(now) {
			next = next.AddDate(0, 0, interval)
		}

		return next.Format(dateFormat), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", repeat)
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
