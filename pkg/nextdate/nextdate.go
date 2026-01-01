package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// 1. Парсим исходную дату
	d, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	// 2. Проверяем, что правило не пустое
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	// 3. Разбираем правило
	// Правило выглядит как "тип значение". Например "d 7"
	parts := strings.Split(repeat, " ")
	ruleType := parts[0]

	switch ruleType {
	case "y":
		return nextDateYear(now, d)
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid d rule format")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days > 400 || days <= 0 {
			return "", fmt.Errorf("invalid days interval")
		}
		return nextDateDays(now, d, days)
	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid w rule format")
		}
		return nextDateWeek(now, d, parts[1])
	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("invalid m rule format")
		}
		// m <дни> [месяцы] -> передаем части после "m"
		return nextDateMonth(now, d, parts[1:])
	default:
		return "", fmt.Errorf("unsupported rule: %s", ruleType)
	}
}

// === Логика для "y" (Год) ===
func nextDateYear(now time.Time, date time.Time) (string, error) {
	// Добавляем по году, пока дата не станет больше now
	for {
		date = date.AddDate(1, 0, 0)
		if date.After(now) {
			break
		}
	}
	return date.Format("20060102"), nil
}

// === Логика для "d" (Дни) ===
func nextDateDays(now time.Time, date time.Time, days int) (string, error) {
	// Добавляем по N дней, пока дата не станет больше now
	for {
		date = date.AddDate(0, 0, days)
		if date.After(now) {
			break
		}
	}
	return date.Format("20060102"), nil
}

// === Логика для "w" (Дни недели) ===
func nextDateWeek(now time.Time, date time.Time, weekdaysStr string) (string, error) {
	var weekdays []int
	for _, s := range strings.Split(weekdaysStr, ",") {
		day, err := strconv.Atoi(s)
		if err != nil || day < 1 || day > 7 {
			return "", fmt.Errorf("invalid weekday: %s", s)
		}
		weekdays = append(weekdays, day)
	}

	// Если задача старая, начинаем перебор с "сегодня", чтобы не крутить цикл зря
	if date.Before(now) {
		date = now
	}

	// Ищем ближайший подходящий день
	for {
		// Сдвигаем на 1 день вперед
		date = date.AddDate(0, 0, 1)

		// Получаем день недели (в Go Sunday=0, Monday=1... но у нас в задании 1=Mon, 7=Sun)
		// Поэтому делаем хитрый маппинг
		currentWd := int(date.Weekday()) // 0=Sun, 1=Mon...
		if currentWd == 0 {
			currentWd = 7
		}

		// Проверяем, подходит ли этот день
		for _, target := range weekdays {
			if currentWd == target {
				return date.Format("20060102"), nil
			}
		}
	}
}

// === Логика для "m" (Дни месяца) ===
// Это самая сложная часть. Мы будем идти вперед по дням и проверять условия.
func nextDateMonth(now time.Time, date time.Time, parts []string) (string, error) {
	// 1. Парсим дни месяца
	var days []int
	for _, s := range strings.Split(parts[0], ",") {
		d, err := strconv.Atoi(s)
		if err != nil || d < -2 || d > 31 || d == 0 {
			return "", fmt.Errorf("invalid month day: %s", s)
		}
		days = append(days, d)
	}

	// 2. Парсим месяцы (если они есть)
	var months []int
	if len(parts) == 2 {
		for _, s := range strings.Split(parts[1], ",") {
			m, err := strconv.Atoi(s)
			if err != nil || m < 1 || m > 12 {
				return "", fmt.Errorf("invalid month: %s", s)
			}
			months = append(months, m)
		}
	}

	// Если задача старая, начинаем перебор с "сегодня"
	if date.Before(now) {
		date = now
	}

	// Ищем дату
	for {
		date = date.AddDate(0, 0, 1)

		// Проверяем месяц (если список месяцев задан)
		if len(months) > 0 {
			currentMonth := int(date.Month())
			found := false
			for _, m := range months {
				if currentMonth == m {
					found = true
					break
				}
			}
			if !found {
				continue // Месяц не подошел, идем к следующему дню
			}
		}

		// Проверяем день
		currentDay := date.Day()
		// Определяем последний день текущего месяца
		// date.AddDate(0, 1, -currentDay) переносит нас в конец текущего месяца
		lastDayOfMonth := date.AddDate(0, 1, -currentDay).Day()

		for _, targetDay := range days {
			if targetDay > 0 {
				// Обычный день (1, 15, 31)
				if currentDay == targetDay {
					return date.Format("20060102"), nil
				}
			} else {
				// Отрицательный день (-1 - последний, -2 - предпоследний)
				// Если targetDay = -1, то lastDayOfMonth - 0
				// Если targetDay = -2, то lastDayOfMonth - 1
				if currentDay == lastDayOfMonth+1+targetDay {
					return date.Format("20060102"), nil
				}
			}
		}
	}
}
