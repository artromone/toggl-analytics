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
