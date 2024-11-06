package config

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"togglparser/internal/types"
)

type CredentialField int

const (
	apiKey CredentialField = iota
	workspaceId
	username
	payPerHour
	clientPay
)

func (d CredentialField) String() string {
	return [...]string{"API_KEY", "WORKSPACE_ID", "USER_NAME", "PAY_PER_HOUR", "CLIENT_PAY"}[d]
}

func GetUserCredentials(prefix string) types.UserCredentials {
	prefix += "_"
	return types.UserCredentials{
		APIKey:      os.Getenv(prefix + apiKey.String()),
		WorkspaceID: os.Getenv(prefix + workspaceId.String()),
		FileName:    os.Getenv(prefix + username.String()),
		PayPerHour:  os.Getenv(prefix + payPerHour.String()),
	}
}

func GetAllUserCredentials() map[string]types.UserCredentials {
	users := make(map[string]types.UserCredentials)

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
			users[userPrefix] = types.UserCredentials{}
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

func CheckCredentials(r types.Requester, apiKey string) error {
	// 1 auth
	url := "https://api.track.toggl.com/api/v9/me"

	resp, err := r.MakeRequest(http.MethodGet, url)
	if err != nil {
		return fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	// 2 check client
	_, exists := os.LookupEnv("CLIENT_PAY")
	if !exists {
		return fmt.Errorf("CLIENT_PAY environment variable is not set")
	}

	return nil
}
