package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Failed to load .env file: %v", err)
	}

	users := GetAllUserCredentials()

	for userPrefix, credentials := range users {
		fmt.Printf("Checking credentials for %s...", userPrefix)
		apiKey := credentials.APIKey

		CheckCredentials(apiKey)
	}
}

func CheckCredentials(apiKey string) {

	method, url := http.MethodGet, "https://api.track.toggl.com/api/v9/me"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		req = nil
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		resp = nil
	}

	if err != nil {
		fmt.Printf(" Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf(" Unexpected status code: %d\n", resp.StatusCode)
	} else {
		fmt.Printf(" Credentials are valid.\n")
	}
}
