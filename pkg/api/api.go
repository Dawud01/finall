package api

import (
	"net/http"
	"time"

	// Импортируем нашу логику вычислений
	"go1f/pkg/nextdate"
)

// Константа для формата даты (чтобы не ошибиться в будущем)
const DateFormat = "20060102"

// Init регистрирует маршруты (ручки) нашего API
// Мы вызываем эту функцию из server.go
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)

	http.HandleFunc("/api/task/done", doneTaskHandler)
}

// nextDateHandler - это наш "Официант"
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем параметры из запроса
	nowStr := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	// 2. Обрабатываем параметр "now"
	if nowStr == "" {
		// Если now не прислали, берем текущее время сервера
		now = time.Now()
	} else {
		// Если прислали, пытаемся распарсить
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			// Если формат кривой — ругаемся (400 Bad Request)
			http.Error(w, "invalid 'now' format: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// 3. Зовем "Шеф-повара" (функцию из прошлого задания)
	// Важно: date и repeat мы передаем как есть, валидация внутри NextDate
	nextDate, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		// Если логика вернула ошибку (например, кривое правило), отдаем её клиенту
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Отдаем результат клиенту
	// w.Write принимает байты, поэтому конвертируем строку
	w.Write([]byte(nextDate))
}
