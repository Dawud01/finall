package api

import (
	"net/http"
	"time"

	"go1f/pkg/db"
	"go1f/pkg/nextdate"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Получаем ID
	id := r.FormValue("id")
	if id == "" {
		sendError(w, "Не указан идентификатор")
		return
	}

	// 2. Ищем задачу, чтобы понять, она одноразовая или цикличная
	task, err := db.GetTask(id)
	if err != nil {
		sendError(w, "Задача не найдена")
		return
	}

	// 3. Логика обработки
	if task.Repeat == "" {
		// Вариант А: Задача одноразовая -> Удаляем
		if err := db.DeleteTask(id); err != nil {
			sendError(w, "Ошибка при удалении задачи: "+err.Error())
			return
		}
	} else {
		// Вариант Б: Задача цикличная -> Переносим дату
		now := time.Now()

		// Рассчитываем следующую дату
		next, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			sendError(w, "Ошибка расчета следующей даты: "+err.Error())
			return
		}

		// Обновляем дату в структуре и сохраняем в БД
		task.Date = next
		if err := db.UpdateTask(task); err != nil {
			sendError(w, "Ошибка при обновлении даты задачи: "+err.Error())
			return
		}
	}

	// 4. Успех
	sendJSON(w, map[string]interface{}{})
}
