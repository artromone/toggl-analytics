package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"togglparser/internal/types"
)

type fetcher struct {
	client *http.Client
	apiKey string
}

func NewFetcher(apiKey string) types.Fetcher {
	return &fetcher{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}

func (f *fetcher) FetchData(url string, v interface{}) error {
	resp, err := f.MakeRequest(http.MethodGet, url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("Failed to fetch data, status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func (f *fetcher) MakeRequest(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(f.apiKey, "api_token")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
