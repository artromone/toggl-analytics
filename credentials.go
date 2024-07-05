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
	PayPerHour  string
}

type CredentialField int

const (
	apiKey CredentialField = iota
	workspaceId
	username
	payPerHour
)

func (d CredentialField) String() string {
	return [...]string{"API_KEY", "WORKSPACE_ID", "USER_NAME", "PAY_PER_HOUR"}[d]
}

func GetUserCredentials(prefix string) UserCredentials {
	prefix += "_"
	return UserCredentials{
		APIKey:      os.Getenv(prefix + apiKey.String()),
		WorkspaceID: os.Getenv(prefix + workspaceId.String()),
		FileName:    os.Getenv(prefix + username.String()),
		PayPerHour:  os.Getenv(prefix + payPerHour.String()),
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
		case apiKey.String():
			user.APIKey = value
		case workspaceId.String():
			user.WorkspaceID = value
		case username.String():
			user.FileName = value
		case payPerHour.String():
			user.PayPerHour = value
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

	return nil
}
