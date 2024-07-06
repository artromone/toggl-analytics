package main

import (
	"fmt"
	"togglparser/internal/api"
	"togglparser/internal/config"
	"togglparser/internal/pdf"
	"togglparser/internal/report"
)

func main() {
	if config.LoadEnv(".env") != nil {
		return
	}

	users := config.GetAllUserCredentials()

	for _, credentials := range users {
		fmt.Printf("Checking credentials for %s...", credentials.FileName)

		err := config.CheckCredentials(credentials.APIKey)
		if err != nil {
			fmt.Printf(" Error checking credentials: %v\n", err)
			fmt.Println("Remove or fix invalid credentials to continue.")
			return
		}
		fmt.Printf(" Credentials are valid.\n")
	}

	fmt.Println()

	table := make(report.Table)

	for _, credentials := range users {
		weekTotal, totalPay, err := api.GetLastWeekTimeEntries(&table, &credentials)
		if err != nil {
			fmt.Printf("Error getting time entry: %v\n", err)
			continue
		}

		hours := weekTotal / 3600
		minutes := (weekTotal % 3600) / 60

		fmt.Printf("%s have worked %d h %d min, need to pay %d.\n", credentials.FileName, hours, minutes, totalPay)
	}
	
    fmt.Println()

	columns, rows, colWidths := pdf.GeneratePdfData(table)
	outputPath := "reports/table.pdf"
	err := pdf.CreateDetailedPdfReport(columns, rows, colWidths, outputPath)
	if err != nil {
		fmt.Printf("Error generating PDF: %v\n", err)
	}
}
