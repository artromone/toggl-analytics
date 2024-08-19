package api

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
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

var projectCache = struct {
	sync.RWMutex
	m map[string]ProjectEntry
}{m: make(map[string]ProjectEntry)}

var clientCache = struct {
	sync.RWMutex
	m map[string]ClientEntry
}{m: make(map[string]ClientEntry)}

var ClientsPay = map[string]float64{}

func Bod(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func GetLastWeekTimeEntries(table *report.Table, credentials *types.UserCredentials) (int, int, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return 0, 0, err
	}

	now := time.Now().In(location)
	thisMonday := Bod(now)
	for thisMonday.Weekday() != time.Monday {
		thisMonday = thisMonday.AddDate(0, 0, -1)
	}

	// fmt.Println(thisMonday.Format("2006-01-02"))
	// thisMonday = thisMonday.AddDate(0, 0, -7) // last week

	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday

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
	clientPayStr := os.Getenv("CLIENT_PAY")

	for _, entry := range timeEntries {
		totalDuration += entry.Duration

		fmt.Println(entry)

		clientId, err := GetProjectClient(entry.Workspace, entry.Project, credentials.APIKey)
		if err != nil {
			fmt.Println("!!! BAD CLIENT_ID !!!")
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

		clientPay, err := strconv.ParseFloat(clientPayStr, 64)
		if err != nil {
			panic("No CLIENT_PAY entry in .env")
		}
		if len(clientName) != 0 {
			ClientsPay[clientName] += clientPay * float64(entry.Duration) / 3600
			// RoundToPrecision(clientsPay[clientName], 0)
		}

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
	key := fmt.Sprintf("%d-%d", workspaceID, projectID)
	projectCache.RLock()
	entry, ok := projectCache.m[key]
	projectCache.RUnlock()

	if ok {
		return entry.Client, nil
	}

	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/projects/%d", workspaceID, projectID)
	var newEntry ProjectEntry
	if err := NewFetcher().FetchData(url, apiKey, &newEntry); err != nil {
		return 0, err
	}

	projectCache.Lock()
	defer projectCache.Unlock()
	projectCache.m[key] = newEntry

	return newEntry.Client, nil
}

func GetClientName(workspaceID, clientID int, apiKey string) (string, error) {
	key := fmt.Sprintf("%d-%d", workspaceID, clientID)
	clientCache.RLock()
	entry, ok := clientCache.m[key]
	clientCache.RUnlock()

	if ok {
		return entry.Name, nil
	}

	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/clients/%d", workspaceID, clientID)
	var newEntry ClientEntry
	if err := NewFetcher().FetchData(url, apiKey, &newEntry); err != nil {
		return "", err
	}
	return newEntry.Name, nil
}
