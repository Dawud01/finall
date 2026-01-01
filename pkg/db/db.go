package db

import (
	"database/sql"
	"log"
	"os"

	// Импортируем драйвер. Подчеркивание _ означает, что мы не используем
	// функции из этого пакета напрямую, но нам нужно, чтобы сработала его функция init(),
	// которая зарегистрирует драйвер "sqlite" в системе.
	_ "modernc.org/sqlite"
)

// DB - глобальная переменная для доступа к базе данных.
// Пишем с большой буквы, чтобы другие пакеты (например, handlers) могли её видеть.
var DB *sql.DB

func Init(dbFile string) error {
	// 1. Проверяем, существует ли файл БД
	install := false
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
		log.Println("Database file not found, creating new one...")
	}

	// 2. Открываем соединение (или создаем файл, если его нет)
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// 3. Проверяем соединение (sql.Open ленивый, он не коннектится сразу)
	if err = DB.Ping(); err != nil {
		return err
	}

	// 4. Если файла не было, создаем структуру (таблицу и индекс)
	if install {
		if err = createSchema(); err != nil {
			return err
		}
		log.Println("Database schema created successfully.")
	}

	return nil
}

// createSchema создает таблицы и индексы
func createSchema() error {
	// SQL-запрос для создания таблицы
	// date CHAR(8) - потому что формат YYYYMMDD всегда 8 символов.
	// repeat VARCHAR(128) - ограничение из задания.
	const schema = `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	
	CREATE INDEX scheduler_date ON scheduler(date);
	`

	// Выполняем SQL запрос
	_, err := DB.Exec(schema)
	return err
}
