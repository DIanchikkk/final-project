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
		return "", fmt.Errorf("ошибка преобразования даты")
	}

	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	if repeat == "y" {
		next := start.AddDate(1, 0, 0)
		for !next.After(now) {
			next = next.AddDate(1, 0, 0)
		}
		return next.Format(dateFormat), nil
	}

	parts := strings.Fields(repeat)

	if parts[0] == "d" {
		if len(parts) != 2 {
			return "", fmt.Errorf("неподдерживаемый формат")
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 || n > 400 {
			return "", fmt.Errorf("неподдерживаемый формат")
		}
		next := start.AddDate(0, 0, n)
		for !next.After(now) {
			next = next.AddDate(0, 0, n)
		}
		return next.Format(dateFormat), nil
	}

	if parts[0] == "w" {
		if len(parts) != 2 {
			return "", fmt.Errorf("неподдерживаемый формат")
		}

		days := map[int]bool{}
		for _, s := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(s)
			if err != nil || d < 1 || d > 7 {
				return "", fmt.Errorf("неподдерживаемый формат")
			}
			days[d] = true
		}

		currentDate := start
		for {
			currentDate = currentDate.AddDate(0, 0, 1)
			if currentDate.After(now) {
				wd := int(currentDate.Weekday())
				if wd == 0 {
					wd = 7
				}
				if days[wd] {
					return currentDate.Format(dateFormat), nil
				}
			}
		}
	}

	if parts[0] == "m" {
		if len(parts) < 2 {
			return "", fmt.Errorf("неподдерживаемый формат")
		}

		days := map[int]bool{}
		for _, s := range strings.Split(parts[1], ",") {
			d, err := strconv.Atoi(s)
			if err != nil || d == 0 || d < -2 || d > 31 {
				return "", fmt.Errorf("неподдерживаемый формат")
			}
			days[d] = true
		}

		months := map[int]bool{}
		if len(parts) == 3 {
			for _, s := range strings.Split(parts[2], ",") {
				m, err := strconv.Atoi(s)
				if err != nil || m < 1 || m > 12 {
					return "", fmt.Errorf("неподдерживаемый формат")
				}
				months[m] = true
			}
		}

		currentDate := start
		for {
			currentDate = currentDate.AddDate(0, 0, 1)
			if !currentDate.After(now) {
				continue
			}

			if len(months) > 0 && !months[int(currentDate.Month())] {
				continue
			}

			day := currentDate.Day()
			last := time.Date(currentDate.Year(), currentDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

			if days[day] ||
				(days[-1] && day == last) ||
				(days[-2] && day == last-1) {
				return currentDate.Format(dateFormat), nil
			}
		}
	}

	return "", fmt.Errorf("неподдерживаемый формат")
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(next))
}
