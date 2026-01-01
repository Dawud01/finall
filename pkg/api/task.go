package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go1f/pkg/db"
	"go1f/pkg/nextdate"
)

// taskHandler распределяет запросы по методам
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		// === ВОТ ЭТО МЫ ДОБАВЛЯЕМ ===
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// === GET: Получение задачи по ID ===
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "Не указан идентификатор")
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, "Задача не найдена")
		return
	}

	sendJSON(w, task)
}

// === POST: Создание новой задачи ===
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// 1. Десериализация
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendError(w, "Ошибка десериализации JSON: "+err.Error())
		return
	}

	// 2. Валидация и расчет дат
	if err := processTaskParams(&task); err != nil {
		sendError(w, err.Error())
		return
	}

	// 3. Сохранение в БД
	id, err := db.AddTask(task)
	if err != nil {
		sendError(w, "Ошибка при сохранении в базу: "+err.Error())
		return
	}

	// 4. Ответ с ID
	resp := map[string]string{
		"id": strconv.FormatInt(id, 10),
	}
	sendJSON(w, resp)
}

// === PUT: Обновление задачи ===
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// 1. Десериализация
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		sendError(w, "Ошибка десериализации JSON: "+err.Error())
		return
	}

	// 2. Проверка ID (обязательно для PUT)
	if task.ID == "" {
		sendError(w, "Не указан идентификатор задачи")
		return
	}

	// 3. Валидация и расчет дат
	if err := processTaskParams(&task); err != nil {
		sendError(w, err.Error())
		return
	}

	// 4. Обновление в БД
	if err := db.UpdateTask(task); err != nil {
		sendError(w, "Ошибка при обновлении задачи: "+err.Error())
		return
	}

	// 5. Успешный пустой ответ
	sendJSON(w, map[string]interface{}{})
}

// === ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ ===

// processTaskParams содержит общую логику проверки Title, Date и Repeat
func processTaskParams(task *db.Task) error {
	// Проверка заголовка
	if task.Title == "" {
		return fmt.Errorf("Не указан заголовок задачи")
	}

	now := time.Now()
	todayStr := now.Format(DateFormat)

	// Если дата пустая -> ставим сегодня
	if task.Date == "" {
		task.Date = todayStr
	}

	// Проверяем формат даты
	parseDate, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("Дата представлена в формате, отличном от %s", DateFormat)
	}

	// Логика переноса дат
	// Сравниваем строки: если дата задачи меньше сегодняшней
	if parseDate.Format(DateFormat) < todayStr {
		if task.Repeat == "" {
			// Если правила нет — переносим на сегодня
			task.Date = todayStr
		} else {
			// Если правило есть — считаем следующую дату от СЕГОДНЯ
			next, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("Правило повторения указано в неправильном формате: %v", err)
			}
			task.Date = next
		}
	} else if task.Repeat != "" {
		// Если дата в будущем, но есть правило — проверяем его валидность
		_, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("Правило повторения указано в неправильном формате: %v", err)
		}
	}

	return nil
}

func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	resp := map[string]string{"error": message}
	json.NewEncoder(w).Encode(resp)
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "Не указан идентификатор")
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		sendError(w, "Ошибка при удалении задачи: "+err.Error())
		return
	}

	sendJSON(w, map[string]interface{}{})
}
