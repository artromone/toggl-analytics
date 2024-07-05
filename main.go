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

		err := CheckCredentials(credentials.APIKey)
		if err != nil {
			fmt.Printf(" Error checking credentials: %v\n", err)
			fmt.Println("Remove or fix invalid credentials to continue.")
			return
		}
		fmt.Printf(" Credentials are valid.\n")
	}

	table := make(Table)

	for _, credentials := range users {
		weekTotal, err := GetLastWeekTimeEntries(&table, &credentials)
		if err != nil {
			fmt.Printf("Error getting time entry: %v\n", err)
			continue
		}

		hours := weekTotal / 3600
		minutes := (weekTotal % 3600) / 60

		fmt.Printf("User have worked %d h %d min.\n", hours, minutes)
	}

    table.PrintTable()

	// columns := []string{"ID", "User", "Time", "Sum", "Project"}
	// rows := [][]string{
	// 	{"1", "John", "25", "New York"},
	// 	{"2", "Emily", "30", "Los Angeles"},
	// 	{"3", "Michael", "28", "Chicago"},
	// 	{"4", "Sophia", "26", "Houston"},
	// 	{"5", "William", "29", "Miami"},
	// }
	//
	// GenerateTablePdf(columns, rows)

	// TODO create pdf in pdf/ dir, gitignore
}
