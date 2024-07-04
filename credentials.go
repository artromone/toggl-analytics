package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type UserCredentials struct {
	APIKey      string
	WorkspaceID string
	FileName    string
}

func GetUserCredentials(prefix string) UserCredentials {
	return UserCredentials{
		APIKey:      os.Getenv(prefix + "_API_KEY"), // TODO one place
		WorkspaceID: os.Getenv(prefix + "_WORKSPACE_ID"),
		FileName:    os.Getenv(prefix + "_FILE_NAME"),
	}
}

func GetAllUserCredentials() map[string]UserCredentials {
	users := make(map[string]UserCredentials)

	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)
		key := kv[0]
		value := kv[1]

		if !strings.HasPrefix(key, "USER") {
			continue
		}

		parts := strings.SplitN(key, "_", 2)
		if len(parts) != 2 {
			continue
		}

		userPrefix := parts[0]
		attribute := parts[1]

		if _, exists := users[userPrefix]; !exists {
			users[userPrefix] = UserCredentials{}
		}

		user := users[userPrefix]
		switch attribute {
		case "API_KEY":
			user.APIKey = value
		case "WORKSPACE_ID":
			user.WorkspaceID = value
		case "FILE_NAME":
			user.FileName = value
		}
		users[userPrefix] = user
	}

	return users
}

func CheckCredentials(apiKey string) error {
	method, url := http.MethodGet, "https://api.track.toggl.com/api/v9/me"

	resp, err := MakeRequest(method, url, apiKey)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	fmt.Printf(" Credentials are valid.\n")
	return nil
}
