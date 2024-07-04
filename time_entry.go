package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TimeEntry struct {
	Duration    int    `json:"duration"`
	Description string `json:"description"`
}

func GetLastWeekTimeEntry(apiKey, dateFormat string) error {
	thisMonday := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday.AddDate(0, 0, -1)

	return GetTimeEntry(apiKey, lastMonday, lastSunday, dateFormat)
}

func GetTimeEntry(apiKey string, startDate, endDate time.Time, dateFormat string) error {
	query := fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/me/time_entries?%s", query)

	resp, err := MakeRequest(http.MethodGet, url, apiKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var timeEntries []TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntries); err != nil {
		return err
	}

	totalDuration := 0
	for _, entry := range timeEntries {
		totalDuration += entry.Duration
		// fmt.Printf("Task description: %s\n", entry.Description)
	}

	hours := totalDuration / 3600
	minutes := (totalDuration % 3600) / 60

	fmt.Printf("User have worked %d h %d min from %s to %s.\n", hours, minutes, startDate.Format(dateFormat), endDate.Format(dateFormat))

	return nil
}
