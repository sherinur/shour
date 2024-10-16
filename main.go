package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var progressToken string
	var platformToken string
	maxHours := 20.0

	green := "\033[32m"
	red := "\033[31m"
	blue := "\033[34m"
	brightYellow := "\033[1;33m"
	lightRed := "\033[1;31m"
	reset := "\033[0m"

	// Получение токенов из аргументов командной строки
	if len(os.Args) >= 3 {
		progressToken = os.Args[1]
		platformToken = os.Args[2]
	} else {
		fmt.Println("Usage: shour [progress token] [platform token] [REQUIRED HOURS]")
		os.Exit(0)
	}

	if len(os.Args) == 4 {
		if os.Args[3] == "30" {
			maxHours = 30.0
		}
	}

	// URL для API
	urlProgress := "https://progress.alem.school/api/v1/user/me"
	urlPlatform := "https://platform.alem.school/api/v1/auth/me"
	urlSlots := "https://platform.alem.school/api/v1/review-slots/"

	// Запрос к URL платформы
	reqPlatform, err := http.NewRequest("GET", urlPlatform, nil)
	if err != nil {
		log.Fatal(err)
	}
	reqPlatform.Header.Set("Authorization", "Bearer "+platformToken)
	reqPlatform.Header.Set("Accept", "application/json")

	client := &http.Client{}
	respPlatform, err := client.Do(reqPlatform)
	if err != nil {
		log.Fatal(err)
	}
	defer respPlatform.Body.Close()

	if respPlatform.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(respPlatform.Body)
		log.Fatalf("Error from platform API: %s", string(body))
	}

	// // Извлечение данных из ответа платформы
	var jsonResponse struct {
		Attrs struct {
			ReviewPoints int `json:"review_points"`
		} `json:"attrs"`
	}

	bodyPlatform, err := ioutil.ReadAll(respPlatform.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(bodyPlatform, &jsonResponse); err != nil {
		log.Fatal(err)
	}

	// Запрос к URL прогресса
	client = &http.Client{}

	reqProgress, err := http.NewRequest("GET", urlProgress, nil)
	if err != nil {
		log.Fatal(err)
	}
	reqProgress.Header.Set("Authorization", "Bearer "+progressToken)
	reqProgress.Header.Set("Accept", "application/json")

	respProgress, err := client.Do(reqProgress)
	if err != nil {
		log.Fatal(err)
	}
	defer respProgress.Body.Close()

	if respProgress.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(respProgress.Body)
		log.Fatalf("Error from progress API: %s", string(body))
	}

	// Извлечение данных из ответа прогресса
	bodyProgress, err := ioutil.ReadAll(respProgress.Body)
	if err != nil {
		log.Fatal(err)
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(bodyProgress, &userInfo); err != nil {
		log.Fatal(err)
	}

	// Форматирование строк для вывода
	lives, livesOk := userInfo["lives"].(float64)
	if !livesOk {
		log.Fatal("Lives field is not present or not a number")
	}

	hearts := ""
	for i := 0; i < int(lives); i++ {
		hearts += "❤️  "
	}

	hours, hoursOk := userInfo["hours"].(float64)
	if !hoursOk {
		log.Fatal("Hours field is not present or not a number")
	}

	isRequirementFulfilled := hours >= maxHours
	hoursRounded := fmt.Sprintf("%.2f", hours)
	filledBlocks := int(hours / maxHours * 10)

	scale := ""
	for i := 0; i < filledBlocks; i++ {
		scale += green + "■" + reset
	}
	for i := filledBlocks; i < 10; i++ {
		scale += red + "■" + reset
	}

	// Вывод информации о студенте с форматированием
	fmt.Println("===== Student Info (by nsheri) =====")
	fmt.Printf("   Name:         %s%s%s\n", brightYellow, userInfo["login"], reset)
	fmt.Printf("   Lives:        %s\n", hearts)
	fmt.Printf("   RP:           %d\n", jsonResponse.Attrs.ReviewPoints)
	fmt.Printf("   Hours:        %s \n", hoursRounded)
	fmt.Println("====================================")
	fmt.Printf("   Scale:        [%s ] \n", scale)

	if isRequirementFulfilled {
		fmt.Println(green + "   ✔ Hours fulfilled!" + reset)
	} else if len(os.Args) == 4 && maxHours == 30.0 {
		fmt.Printf(lightRed+"   ✖ Hours left: %.2f%s\n", 30-hours, reset)
	} else {
		fmt.Printf(lightRed+"   ✖ Hours left: %.2f%s\n", 20-hours, reset)
	}
	// Запрос к API слотов
	reqSlots, err := http.NewRequest("GET", urlSlots, nil)
	if err != nil {
		log.Fatal(err)
	}
	reqSlots.Header.Set("Authorization", "Bearer "+platformToken)
	reqSlots.Header.Set("Accept", "application/json")

	client = &http.Client{}
	respSlots, err := client.Do(reqSlots)
	if err != nil {
		log.Fatal(err)
	}
	defer respSlots.Body.Close()

	if respSlots.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(respSlots.Body)
		log.Fatalf("Error from review slots API: %s", string(body))
	}

	// Извлечение данных из ответа
	bodySlots, err := ioutil.ReadAll(respSlots.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseSlots struct {
		Slots []struct {
			StartAt string `json:"start_at"`
			EndAt   string `json:"end_at"`
			Reviews []struct {
				Login string `json:"login"`
			} `json:"reviews"`
		} `json:"slots"`
	}

	if err := json.Unmarshal(bodySlots, &responseSlots); err != nil {
		log.Fatal(err)
	}

	// Текущая дата для сравнения
	now := time.Now()

	// Проверка наличия ревью с текущей даты и до будущих дней
	reviewsExist := false
	for _, slot := range responseSlots.Slots {
		// Парсинг даты начала слота
		startTime, err := time.Parse(time.RFC3339, slot.StartAt)
		if err != nil {
			log.Fatalf("Error parsing time: %v", err)
		}

		// Если слот начинается после текущей даты
		if startTime.After(now) {
			if len(slot.Reviews) > 0 {
				reviewsExist = true
				break
			}
		}
	}

	// Вывод результата
	if reviewsExist {
		fmt.Println(blue + "   ✔ You have upcoming reviews." + reset)
	} else {
		fmt.Println("   ✖ No upcoming reviews.")
	}

	fmt.Println("====================================")
}
