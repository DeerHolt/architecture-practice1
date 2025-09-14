package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	myData := strings.Split(data, ",")
	if len(myData) != 3 {
		return 0, "", 0, fmt.Errorf("data is invalid")
	}

	steps, err := strconv.Atoi(myData[0])
	if err != nil {
		log.Println(err)
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("steps must be greater than 0")
	}

	duration, err := time.ParseDuration(myData[2])
	if err != nil {
		log.Println(err)
		return 0, "", 0, err
	}
	if duration.Nanoseconds() <= 0 {
		return 0, "", 0, fmt.Errorf("duration is not positive")
	}

	return steps, myData[1], duration, nil
}

func distance(steps int, height float64) float64 {
	if height <= 0 {
		return 0.0
	}
	stepLength := height * stepLengthCoefficient
	if steps <= 0 {
		return 0.0
	}

	totalDistance := float64(steps) * stepLength
	totalDistance = totalDistance / float64(mInKm)
	return totalDistance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Nanoseconds() <= 0 {
		return 0.0
	}
	totalDistance := distance(steps, height)
	if totalDistance <= 0 {
		return 0.0
	}

	result := totalDistance / duration.Hours()
	return result
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	var calories float64
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	totalDistance := distance(steps, height)
	averageSpeed := meanSpeed(steps, height, duration)

	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), totalDistance, averageSpeed, calories)
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps must be greater than 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be greater than 0")
	}
	if duration.Nanoseconds() <= 0 {
		return 0, fmt.Errorf("duration is not positive")
	}

	averageSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	resultCalories := (averageSpeed * weight * durationInMinutes) / minInH
	return resultCalories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("steps must be greater than 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("weight must be greater than 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("height must be greater than 0")
	}
	if duration.Nanoseconds() <= 0 {
		return 0, fmt.Errorf("duration is not positive")
	}

	averageSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	resultCalories := (averageSpeed * weight * durationInMinutes) / minInH * walkingCaloriesCoefficient
	return resultCalories, nil
}
