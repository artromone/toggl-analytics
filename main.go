package main

import (
	"fmt"
	"strconv"
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

	// table.PrintTable()

	columns := []string{"ID", "User", "Duration", "Sum", "Client", "Task", "Vikunja link"}

	rows := [][]string{}
	for id, row := range table {
		user, ok := row["User"].(string)
		if !ok {
			continue
		}
		duration, ok := row["Duration"].(int)
		if !ok {
			continue
		}
		sum, ok := row["Sum"].(float64)
		if !ok || sum == 0 {
			continue
		}
		project, ok := row["Client"].(string)
		if !ok {
			continue
		}
		task, ok := row["Task"].(string)
		if !ok {
			continue
		}
		// TODO link together
		taskTrackerId, ok := row["Vikunja link"].(int)
		if !ok {
			continue
		}

		newRow := []string{
			fmt.Sprintf("%d", id),
			user,
			DurationToHHMMSS(duration),
			fmt.Sprintf("%.2f", sum),
			project,
			task,
			strconv.Itoa(taskTrackerId),
		}
		rows = append(rows, newRow)
	}

	colWidths := map[int]float64{
		0: 5.0,
		1: 20.0,
		2: 15.0,
		3: 15.0,
		4: 25.0,
		5: 70.0,
		6: 20.0,
	}

	GenerateTablePdf(columns, rows, colWidths)
}

func DurationToHHMMSS(durationInSeconds int) string {
	hours := durationInSeconds / 3600
	minutes := (durationInSeconds % 3600) / 60
	seconds := durationInSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
