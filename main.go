package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"time"
)

type TimeEntry struct {
	Duration int `json:"duration"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file: %v", err)
	}

	users := GetAllUserCredentials()

	for _, credentials := range users {
		fmt.Printf("Checking credentials for %s...", credentials.FileName)
		apiKey := credentials.APIKey

		CheckCredentials(apiKey)

		GetTimeEntry(apiKey)
	}
}

func GetTimeEntry(apiKey string) {

	thisMonday := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday.AddDate(0, 0, -1)

	query := fmt.Sprintf("start_date=%s&end_date=%s", lastMonday.Format("2006-01-02"), lastSunday.Format("2006-01-02"))
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/me/time_entries?%s", query)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Print(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()

	var timeEntries []TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntries); err != nil {
		fmt.Print(err)
	}

	totalDuration := 0
	for _, entry := range timeEntries {
		totalDuration += entry.Duration
	}

	hours := totalDuration / 3600
	minutes := (totalDuration % 3600) / 60

	fmt.Printf("Пользователь отработал %d часов и %d минут на прошлой неделе.\n", hours, minutes)
}

func CheckCredentials(apiKey string) {

	method, url := http.MethodGet, "https://api.track.toggl.com/api/v9/me"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		req = nil
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resp = nil
	}

	if err != nil {
		fmt.Printf(" Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf(" Unexpected status code: %d\n", resp.StatusCode)
	} else {
		fmt.Printf(" Credentials are valid.\n")
	}
}
