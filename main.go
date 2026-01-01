package main

import (
	"log"

	// Импортируем наш новый пакет
	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	// 1. Инициализируем базу данных
	// Файл создастся в корне проекта
	err := db.Init("scheduler.db")
	if err != nil {
		// Если база не открылась — нет смысла запускать сервер. Падаем.
		log.Fatal(err)
	}

	// Хорошая практика: закрыть соединение, когда main завершится.
	// defer откладывает выполнение этой строки до конца работы функции main.
	defer db.DB.Close()

	log.Println("Database initialized.")

	// 2. Запускаем веб-сервер
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
