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
		fmt.Printf("Credentials for %s: %+v\n", userPrefix, credentials)
	}

	dummyUser := users["USER1"]

	apiKey := dummyUser.APIKey

	req, err := http.NewRequest(http.MethodGet,
		"https://api.track.toggl.com/api/v9/me", nil)
	if err != nil {
		print(err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.SetBasicAuth(apiKey, "api_token")

    client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
	}

	fmt.Println("Response Status:", resp.Status)
}
