package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TimeEntry struct {
	Duration  int `json:"duration"`
	Project   int `json:"project_id"`
	Workspace int `json:"workspace_id"`
}

type ProjectEntry struct {
	Name   string `json:"name"`
	Client int    `json:"client_id"`
}

type ClientEntry struct {
	Name string `json:"name"`
}

func GetLastWeekTimeEntries(table *Table, credentials *UserCredentials) (int, error) {
	thisMonday := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday.AddDate(0, 0, -1)

	return GetTimeEntries(table, credentials, lastMonday, lastSunday)
}

func GetTimeEntries(table *Table, credentials *UserCredentials, startDate, endDate time.Time) (int, error) {
	apiDateFormat := "2006-01-02"
	query := fmt.Sprintf("start_date=%s&end_date=%s&meta=1&include_sharing=1", startDate.Format(apiDateFormat), endDate.Format(apiDateFormat))
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

		clientId, err := GetProjectClient(entry.Workspace, entry.Project, credentials.APIKey)
		if err != nil {
			continue
		}

		clientName, err := GetClientName(entry.Workspace, clientId, credentials.APIKey)
		if err != nil {
			continue
		}

		table.AddRow(credentials.FileName, startDate, float64(entry.Duration)*250/3600, clientName)
	}

	return totalDuration, nil

	// TODO Create data struct with all fields, fill them and keep/send
}

func GetProjectClient(workspaceID, projectID int, apiKey string) (int, error) {
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/projects/%d", workspaceID, projectID)
	resp, err := MakeRequest(http.MethodGet, url, apiKey)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var entry ProjectEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return 0, err
	}

	return entry.Client, nil
}

func GetClientName(workspaceID, clientID int, apiKey string) (string, error) {
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/clients/%d", workspaceID, clientID)
	resp, err := MakeRequest(http.MethodGet, url, apiKey)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var clientEntry ClientEntry
	if err := json.NewDecoder(resp.Body).Decode(&clientEntry); err != nil {
		return "", err
	}

    print(clientEntry.Name)

	return clientEntry.Name, nil
}
