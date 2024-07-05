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

func GetLastWeekTimeEntries(table *Table, credentials *UserCredentials) (int, error) {
	thisMonday := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday.AddDate(0, 0, -1)

	return GetTimeEntries(table, credentials, lastMonday, lastSunday)
}

func GetTimeEntries(table *Table, credentials *UserCredentials, startDate, endDate time.Time) (int, error) {
	apiDateFormat := "2006-01-02"
	query := fmt.Sprintf("start_date=%s&end_date=%s", startDate.Format(apiDateFormat), endDate.Format(apiDateFormat))
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/me/time_entries?%s", query)

	resp, err := MakeRequest(http.MethodGet, url, credentials.APIKey)
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

		// dateRange := startDate.Format(time.RFC3339) + " - " + endDate.Format(time.RFC3339)

		table.AddRow(credentials.FileName, startDate, float64(entry.Duration)*250/3600, entry.Description)
	}

	return totalDuration, nil

	// TODO Create data struct with all fields, fill them and keep/send
}
