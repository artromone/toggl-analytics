package main

import (
	"fmt"
)

func main() {
	if LoadEnv(".env") != nil {
		return
	}

	users := GetAllUserCredentials()

	for _, credentials := range users {
		fmt.Printf("Checking credentials for %s...\n", credentials.FileName)
		apiKey := credentials.APIKey

		err := CheckCredentials(apiKey)
		if err != nil {
			fmt.Printf("Error checking credentials: %v\n", err)
			continue
		}

		err = GetTimeEntry(apiKey)
		if err != nil {
			fmt.Printf("Error getting time entry: %v\n", err)
			continue
		}
	}
}
