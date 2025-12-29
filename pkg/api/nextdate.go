package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

var weekdays [8]bool

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования даты: %v", err)
	}

	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	if repeat == "y" {
		for !start.After(now) {
			start = start.AddDate(1, 0, 0)
		}
		return start.Format(dateFormat), nil
	}

	parts := strings.Split(repeat, " ")
	if parts[0] != "d" || len(parts) < 2 {
		return "", fmt.Errorf("неподдерживаемый формат: %s", repeat)
	}

	twoParts := strings.Split(parts[1], ",")
	for _, part := range twoParts {
		dayNum, err := strconv.Atoi(part)
		if err != nil {
			return "", fmt.Errorf("не удалось преобразовать день недели в число: %v", err)
		}

		if dayNum < 1 || dayNum > 7 {
			return "", fmt.Errorf("день недели должен быть от 1 до 7, получено: %d", dayNum)
		}

		weekdays[dayNum] = true
	}

	for !start.After(now) {
		start = start.AddDate(0, 0, 1)
		weekday := int(start.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		for !weekdays[weekday] {
			start = start.AddDate(0, 0, 1)
			weekday = int(start.Weekday())
			if weekday == 0 {
				weekday = 7
			}
		}
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("не удалось преобразовать интервал в число: %v", err)
	}

	if interval < 1 || interval > 400 {
		return "", fmt.Errorf("интервал должен быть от 1 до 400 дней, получено: %d", interval)
	}

	for !start.After(now) {
		start = start.AddDate(0, 0, interval)
	}

	return start.Format(dateFormat), nil
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
