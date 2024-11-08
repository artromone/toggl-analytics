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

type ProjectEntry = types.ProjectEntry
type ClientEntry = types.ClientEntry
type TimeEntry = types.TimeEntry

var projectCache sync.Map
var clientCache sync.Map
var ClientsPay = map[string]float64{}
// type AppContext struct {
//     ProjectCache sync.Map
//     ClientCache  sync.Map
//     ClientsPay   map[string]float64
// }

const ApiDateFormat = "2006-01-02"

func Bod(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}

func StartOfWeek(t time.Time) time.Time {
	offset := (int(t.Weekday()) - int(time.Monday) + 7) % 7
	return t.AddDate(0, 0, -offset).Truncate(24 * time.Hour)
}

func GetLastWeekTimeEntries(table *report.Table, credentials *types.UserCredentials) (int, int, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return 0, 0, err
	}

	now := time.Now().In(location)
	thisMonday := StartOfWeek(now).AddDate(0, 0, -7) // last week

	lastMonday := thisMonday.AddDate(0, 0, -7)
	lastSunday := thisMonday

	timeEntries, err := GetTimeEntries(credentials, lastMonday, lastSunday)
	if err != nil {
		return 0, 0, err
	}

	return ProcessTimeEntries(table, credentials, timeEntries)
}

func GetTimeEntries(credentials *types.UserCredentials, startDate, endDate time.Time) ([]TimeEntry, error) {
	query := fmt.Sprintf("start_date=%s&end_date=%s&meta=1&include_sharing=1", startDate.Format(ApiDateFormat), endDate.Format(ApiDateFormat))
	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/me/time_entries?%s", query)

	var timeEntries []TimeEntry
	if err := NewFetcher(credentials.APIKey).FetchData(url, &timeEntries); err != nil {
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

		clientId, err := GetProjectClient(entry.Workspace, entry.Project, credentials.APIKey)
		if err != nil {
			fmt.Printf("Error getting clientId for project %d: %v", entry.Project, err)
		}
		clientName, err := GetClientName(entry.Workspace, clientId, credentials.APIKey)
		if err != nil {
			clientName = ""
		}

		pay, err := strconv.ParseFloat(credentials.PayPerHour, 64)
		if err != nil {
			continue
		}
		pay *= DurationToHours(entry.Duration)

		totalPay += int(pay)

		clientPay, err := strconv.ParseFloat(clientPayStr, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("CLIENT_PAY must be a valid number: %v", err)
		}
		if len(clientName) != 0 {
			ClientsPay[clientName] += clientPay * DurationToHours(entry.Duration)
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

func DurationToHours(duration int) float64 {
	return RoundToPrecision(float64(duration)/3600, 2)
}

func RoundToPrecision(value float64, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	return math.Round(value*multiplier) / multiplier
}

// func GetFromCacheOrFetch[T any](cache *sync.Map, key string, fetchFunc func() (T, error)) (T, error) {
//     if value, ok := cache.Load(key); ok {
//         return value.(T), nil
//     }
//
//     result, err := fetchFunc()
//     if err != nil {
//         return *new(T), err
//     }
//
//     cache.Store(key, result)
//     return result, nil
// }

func GetProjectClient(workspaceID, projectID int, apiKey string) (int, error) {
	key := fmt.Sprintf("%d-%d", workspaceID, projectID)
	if value, ok := projectCache.Load(key); ok {
		entry := value.(ProjectEntry)
		return entry.Client, nil
	}

	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/projects/%d", workspaceID, projectID)
	var newEntry ProjectEntry
	if err := NewFetcher(apiKey).FetchData(url, &newEntry); err != nil {
		return 0, err
	}

	projectCache.Store(key, newEntry)

	return newEntry.Client, nil
}

func GetClientName(workspaceID, clientID int, apiKey string) (string, error) {
	key := fmt.Sprintf("%d-%d", workspaceID, clientID)
	if value, ok := clientCache.Load(key); ok {
		entry := value.(ClientEntry)
		return entry.Name, nil
	}

	url := fmt.Sprintf("https://api.track.toggl.com/api/v9/workspaces/%d/clients/%d", workspaceID, clientID)
	var newEntry ClientEntry
	if err := NewFetcher(apiKey).FetchData(url, &newEntry); err != nil {
		return "", err
	}

	clientCache.Store(key, newEntry)

	return newEntry.Name, nil
}
