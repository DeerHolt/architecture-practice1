package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage проверяет, что полученные данные корректны
func parsePackage(data string) (int, time.Duration, error) {
	myData := strings.Split(data, ",")
	if len(myData) != 2 {
		return 0, 0, fmt.Errorf("data is invalid")
	}

	steps, err := strconv.Atoi(myData[0])
	if err != nil {
		log.Println(err)
		return 0, 0, fmt.Errorf("steps is not integer")
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("steps must be greater than 0")
	}

	duration, err := time.ParseDuration(myData[1])
	if err != nil {
		log.Println(err)
		return 0, 0, fmt.Errorf("duration does not meet time format")
	}
	if duration.Nanoseconds() <= 0 {
		return 0, 0, fmt.Errorf("duration is not positive")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if weight <= 0 {
		fmt.Println("weight must be greater than 0")
		return ""
	}
	if height <= 0 {
		fmt.Println("height must be greater than 0")
		return ""
	}

	distance := float64(steps) * stepLength
	distance = distance / float64(mInKm)

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
