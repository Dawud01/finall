package server

import (
	"log"
	"net/http"
	"os"

	// Добавляем импорт нашего нового пакета API
	"go1f/pkg/api"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// === ВОТ ЭТУ СТРОКУ НУЖНО ДОБАВИТЬ ===
	// Мы регистрируем API обработчики ДО запуска сервера
	api.Init()
	// =====================================

	dir := http.Dir("./web")
	handler := http.FileServer(dir)
	http.Handle("/", handler)

	log.Printf("Server is running on port %s...", port)
	return http.ListenAndServe(":"+port, nil)
}
