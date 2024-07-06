package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"togglparser/internal/types"
)

type fetcher struct{}

func (f *fetcher) FetchData(url, apiKey string, v interface{}) error {
	resp, err := MakeRequest(http.MethodGet, url, apiKey)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to fetch data, status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

func NewFetcher() types.Fetcher {
	return &fetcher{}
}

func MakeRequest(method, url, apiKey string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
