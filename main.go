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
		fmt.Printf("Checking credentials for %s...", credentials.FileName)
		apiKey := credentials.APIKey

		err := CheckCredentials(apiKey)
		if err != nil {
			fmt.Printf(" Error checking credentials: %v\n", err)
			continue
		}

		err = GetLastWeekTimeEntry(apiKey, "02.01.2006")
		if err != nil {
			fmt.Printf("Error getting time entry: %v\n", err)
			continue
		}
	}
}
