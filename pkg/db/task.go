package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const dateLayout = "02012006"
const dbLayout = "20060102"

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	if limit < 1 {
		limit = 50
	}

	var tasks []*Task
	var rows *sql.Rows
	var err error

	search = strings.TrimSpace(search)

	if search == "" {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)
	} else {
		t, errDate := time.Parse(dateLayout, search)
		if errDate == nil {
			dateStr := t.Format(dbLayout)
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`
			rows, err = DB.Query(query, dateStr, limit)
		} else {
			searchPattern := "%" + search + "%"
			query := `SELECT id, date, title, comment, repeat 
			          FROM scheduler 
			          WHERE title LIKE ? OR comment LIKE ? 
			          ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, searchPattern, searchPattern, limit)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к базе: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		task := &Task{}
		if err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("ошибка чтения строки: %v", err)
		}
		task.ID = strconv.FormatInt(id, 10)
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации строк: %v", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	task := &Task{}
	err := DB.QueryRow(query, id).
		Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, fmt.Errorf("задача с id=%s не найдена", id)
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler 
	          SET date = ?, title = ?, comment = ?, repeat = ? 
	          WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %v", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %v", err)
	}
	if affected == 0 {
		return fmt.Errorf("задача с id=%s не найдена", task.ID)
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("задача с id=%s не найдена", id)
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("задача с id=%s не найдена", id)
	}
	return nil
}
