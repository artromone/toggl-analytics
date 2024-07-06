package types

type UserCredentials struct {
	APIKey      string
	WorkspaceID string
	FileName    string
	PayPerHour  string
}

type Fetcher interface {
	FetchData(url, apiKey string, v interface{}) error
}
