package daysteps

import (
	"time"
	"strings"
	"fmt"
	"strconv"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
)

const (
	// Длина одного шага в метрах
	stepLength float64 = 0.65
	// Количество метров в одном километре
	mInKm float64 = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделяем данные
	splittedData := strings.Split(data, ",")

	// Проверяем что получилось 2 элемента
	if len(splittedData) != 2 {
		return 0, 0, fmt.Errorf("bad data: count of args %d != %d", len(splittedData), 2)
	}

	// Парсим шаги и проверяем чтобы они были > 0
	steps, err := strconv.Atoi(splittedData[0])
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("bad data: steps")
	}

	// Парсим продолжительность и проверяем чтобы она была > 0
	timeDur, err := time.ParseDuration(splittedData[1])
	if err != nil || timeDur <= 0 {
		return 0, 0, fmt.Errorf("bad data: time")
	}

	return steps, timeDur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о шагах и продолжительности от parsePakage()
	// И обрабатываем возможные ошибки
	steps, timeDur, err := parsePackage(data)
	if err != nil || steps == 0 {
		log.Println(err)
		return ""
	}

	// Вычисляем и переводим дистанцию в километры
	tripLenght := (float64(steps) * stepLength) / mInKm

	// Вычиcляем калории с помощью WalkingSpentCalories()
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, timeDur)
	if err != nil || calories == 0.0{
		log.Println(err)
		return ""
	}

	// Формируем и возвращаем результирующую строку
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, tripLenght, calories)
}
