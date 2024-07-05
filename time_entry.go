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

func GetLastWeekTimeEntries(apiKey string) (int, error) {
	thisMonday := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday.AddDate(0, 0, -1)

	return GetTimeEntries(apiKey, lastMonday, lastSunday)
}

func GetTimeEntries(apiKey string, startDate, endDate time.Time) (int, error) {
    // TODO one format
	query := fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/me/time_entries?%s", query)

	resp, err := MakeRequest(http.MethodGet, url, apiKey)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var timeEntries []TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntries); err != nil {
		return 0, err
	}

	totalDuration := 0
	for _, entry := range timeEntries {
		totalDuration += entry.Duration
	}

    return totalDuration, nil

    // TODO Create data struct with all fields, fill them and keep/send
}
