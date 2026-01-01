package api

import (
	"net/http"

	"go1f/pkg/db"
)

// tasksHandler обрабатывает запросы на получение списка задач
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Разрешаем только GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Идем в базу за задачами (берем 50 штук, как рекомендовано)
	tasks, err := db.NextTasks(50)
	if err != nil {
		sendError(w, "Ошибка получения задач: "+err.Error())
		return
	}

	// 3. Формируем ответ
	// Нам нужно вернуть объект { "tasks": [...] }
	resp := map[string]interface{}{
		"tasks": tasks,
	}

	// 4. Отправляем JSON
	sendJSON(w, resp)
}
