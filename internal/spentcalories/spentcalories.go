package spentcalories

import (
	"time"
	"strings"
	"fmt"
	"strconv"
	"log"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm float64              = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Парсим строку и проверяем количество полученных эл-ов результирующего массива
	splittedData := strings.Split(data, ",")
	if len(splittedData) != 3 {
		return 0, "", 0, fmt.Errorf("bad data: count of args %d != %d", len(splittedData), 3)
	}

	// Парсим шаги и обрабатываем возможные ошибки
	steps, err := strconv.Atoi(splittedData[0])
	if err != nil || steps <= 0{
		return 0, "", 0, fmt.Errorf("bad data: steps")
	}

	// Парсим продолжительность и обрабатываем возможные ошибки
	timeDur, err := time.ParseDuration(splittedData[2])
	if err != nil || timeDur <= 0 {
		return 0, "", 0, fmt.Errorf("bad data: time")
	}

	// Возвращаем полученные данные
	return steps, splittedData[1], timeDur, nil
}

func distance(steps int, height float64) float64 {
	if steps < 0 || height < 0 {
		return 0
	}

	stepLength := height * stepLengthCoefficient
	return stepLength * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps < 0 || height < 0 || duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получаем данные
	steps, tranningType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	dist := distance(steps, height)
	meanV := meanSpeed(steps, height, duration)
	// Проверяем тип тренировки и обрабатываем возможные ошибки
	switch tranningType {
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
		}
		outStr := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		tranningType, duration.Hours(), dist, meanV, calories)
		return outStr, nil

	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
		}
		outStr := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		tranningType, duration.Hours(), dist, meanV, calories)
		return outStr, nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных данных
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("bad data")
	}
	
	// Подсчет средней скорости и затраченных калорий
	meanV := meanSpeed(steps, height, duration)
	return (weight * meanV * duration.Minutes()) / float64(minInH), nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных данных
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("bad data")
	}
	
	// Подсчет средней скорости и затраченных калорий
	meanV := meanSpeed(steps, height, duration)
	return (weight * meanV * duration.Minutes()) / float64(minInH) * float64(walkingCaloriesCoefficient), nil
}
