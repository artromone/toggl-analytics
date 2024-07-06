package api

import (
	"fmt"
	"math"
	"strconv"
	"time"
	"togglparser/internal/report"
	"togglparser/internal/types"
)

type TimeEntry struct {
	Duration      int    `json:"duration"`
	Project       int    `json:"project_id"`
	Workspace     int    `json:"workspace_id"`
	Task          string `json:"description"`
	TaskTrackerID []int  `json:"tag_ids"`
}

type ProjectEntry struct {
	Name   string `json:"name"`
	Client int    `json:"client_id"`
}

type ClientEntry struct {
	Name string `json:"name"`
}

func GetLastWeekTimeEntries(table *report.Table, credentials *types.UserCredentials) (int, int, error) {
	thisMonday := time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1)
	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday.AddDate(0, 0, -1)

	timeEntries, err := GetTimeEntries(credentials, lastMonday, lastSunday)
	if err != nil {
		return 0, 0, err
	}

	return ProcessTimeEntries(table, credentials, timeEntries)
}

func GetTimeEntries(credentials *types.UserCredentials, startDate, endDate time.Time) ([]TimeEntry, error) {
	apiDateFormat := "2006-01-02"
	query := fmt.Sprintf("start_date=%s&end_date=%s&meta=1&include_sharing=1", startDate.Format(apiDateFormat), endDate.Format(apiDateFormat))
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/me/time_entries?%s", query)

	var timeEntries []TimeEntry
	if err := NewFetcher().FetchData(url, credentials.APIKey, &timeEntries); err != nil {
		return nil, err
	}

	return timeEntries, nil
}

func ProcessTimeEntries(table *report.Table, credentials *types.UserCredentials, timeEntries []TimeEntry) (int, int, error) {
	tasks := make(map[string]int)
	totalDuration := 0

    totalPay := 0

	for _, entry := range timeEntries {
		totalDuration += entry.Duration

		clientId, err := GetProjectClient(entry.Workspace, entry.Project, credentials.APIKey)
		if err != nil {
			continue
		}

		clientName, err := GetClientName(entry.Workspace, clientId, credentials.APIKey)
		if err != nil {
			clientName = ""
		}

		pay, err := strconv.ParseFloat(credentials.PayPerHour, 64)
		if err != nil {
			continue
		}
		pay *= float64(entry.Duration) / 3600
		pay = RoundToPrecision(pay, 0)

        totalPay += int(pay)

		if rowId, exists := tasks[entry.Task]; exists {
			table.UpdateRow(rowId, "Sum", table.Get(rowId, "Sum").(float64)+pay)
			table.UpdateRow(rowId, "Duration", table.Get(rowId, "Duration").(int)+entry.Duration)
		} else {
			taskTrackerId := -1
			if len(entry.TaskTrackerID) != 0 {
				taskTrackerId = entry.TaskTrackerID[0]
			}

			rowId := table.AddRow(
				credentials.FileName,
				entry.Duration,
				pay,
				clientName,
				entry.Task,
				taskTrackerId,
			)
			tasks[entry.Task] = rowId
		}
	}

	return totalDuration, totalPay, nil
}

func RoundToPrecision(value float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(value*multiplier) / multiplier
}

func GetProjectClient(workspaceID, projectID int, apiKey string) (int, error) {
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/projects/%d", workspaceID, projectID)
	var entry ProjectEntry
	if err := NewFetcher().FetchData(url, apiKey, &entry); err != nil {
		return 0, err
	}
	return entry.Client, nil
}

func GetClientName(workspaceID, clientID int, apiKey string) (string, error) {
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/clients/%d", workspaceID, clientID)
	var clientEntry ClientEntry
	if err := NewFetcher().FetchData(url, apiKey, &clientEntry); err != nil {
		return "", err
	}
	return clientEntry.Name, nil
}
