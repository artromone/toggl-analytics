package main

import (
	"fmt"
	"time"
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

	columns := []string{"ID", "User", "Time", "Sum", "Project"}

    rows := [][]string{}
	for id, row := range table {
		user, ok := row["User"].(string)
		if !ok {
			continue
		}
		timeVal, ok := row["Time"].(time.Time)
		if !ok {
			continue
		}
		sum, ok := row["Sum"].(float64)
		if !ok {
			continue
		}
		project, ok := row["Project"].(string)
		if !ok {
			continue
		}

		newRow := []string{fmt.Sprintf("%d", id), user, timeVal.Format(time.RFC3339), fmt.Sprintf("%.2f", sum), project}
		rows = append(rows, newRow)
	}

	GenerateTablePdf(columns, rows)

	// TODO create pdf in pdf/ dir, gitignore
}
