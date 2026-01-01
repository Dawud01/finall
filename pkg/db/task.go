package db

import (
	"fmt"
)

// Task — структура задачи
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в базу данных
func AddTask(t Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// NextTasks возвращает ближайшие задачи (сортировка по дате)
func NextTasks(limit int) ([]Task, error) {
	var tasks []Task
	tasks = []Task{} // Инициализируем, чтобы вернулся [], а не null

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTask возвращает задачу по ID (Нужна для редактирования)
func GetTask(id string) (Task, error) {
	var t Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	// Scan сразу заполняет структуру
	err := DB.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return Task{}, err
	}
	return t, nil
}

// UpdateTask обновляет задачу (Нужна для сохранения изменений)
func UpdateTask(t Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`

	res, err := DB.Exec(query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return err
	}

	// Проверяем, нашлась ли такая задача
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

// DeleteTask удаляет задачу по ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`

	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	// Проверяем, была ли удалена строка (если ID не найден, вернем ошибку)
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
