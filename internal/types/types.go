package types

import (
	"net/http"
)

type UserCredentials struct {
	APIKey      string
	WorkspaceID string
	FileName    string
	PayPerHour  string
}

type Fetcher interface {
	FetchData(url string, v interface{}) error
	MakeRequest(method, url string) (*http.Response, error)
}

type Requester interface {
	MakeRequest(method, url string) (*http.Response, error)
}

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
