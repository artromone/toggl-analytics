package pdf

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"os"
	"sort"
	"strconv"
	"togglparser/internal/report"
)

func CreatePdfReport(columns []string, rows [][]string, colWidths map[int]float64, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "font/")

	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 7)

	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	rowHeight := 4.0

	for i, col := range columns {
		width, exists := colWidths[i]
		if !exists {
			width = 20.0
		}
		pdf.CellFormat(width, rowHeight, tr(col), "1", 0, "", false, 0, "")
	}
	pdf.Ln(rowHeight)

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][0] < rows[j][0] // Sort by ID
	})

	serverAddress := os.Getenv("SERVER_ADDRESS")

	for _, row := range rows {
		for i, col := range row {
			width, exists := colWidths[i]
			if !exists {
				width = 20.0
			}

			linkStr := ""
			value := tr(col)

			isLastColumn := i == len(row)-1
			if isLastColumn {
				i, err := strconv.Atoi(col)
				if err != nil {
					continue
				}

				linkStr = fmt.Sprintf("http://%s/tasks/%d", serverAddress, i)
				value = tr(fmt.Sprintf("Перейти к %d", i))
			}

			pdf.CellFormat(width, rowHeight, value, "1", 0, "", false, 0, linkStr)
		}
		pdf.Ln(rowHeight)
	}

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return fmt.Errorf("Error creating PDF: %v", err)
	}

	fmt.Println("PDF created successfully")
	return nil
}

func GeneratePdfData(table report.Table) (columns []string, rows [][]string, colWidths map[int]float64) {
	columns = []string{"ID", report.UserKey, report.DurationKey, report.SumKey, report.ClientKey, report.TaskKey, report.TaskTrackerKey}

	for id, row := range table {
		user, ok := row[report.UserKey].(string)
		if !ok {
			continue
		}
		duration, ok := row[report.DurationKey].(int)
		if !ok {
			continue
		}
		sum, ok := row[report.SumKey].(float64)
		if !ok || sum == 0 {
			continue
		}
		client, ok := row[report.ClientKey].(string)
		if !ok {
			continue
		}
		task, ok := row[report.TaskKey].(string)
		if !ok {
			continue
		}
		taskTrackerId, ok := row[report.TaskTrackerKey].(int)
		if !ok {
			continue
		}

		newRow := []string{
			fmt.Sprintf("%d", id),
			user,
			DurationToHHMMSS(duration),
			fmt.Sprintf("%.2f", sum),
			client,
			task,
			strconv.Itoa(taskTrackerId),
		}
		rows = append(rows, newRow)
	}

	colWidths = map[int]float64{
		0: 5.0,
		1: 20.0,
		2: 15.0,
		3: 15.0,
		4: 25.0,
		5: 70.0,
		6: 25.0,
	}

	return columns, rows, colWidths
}

func DurationToHHMMSS(durationInSeconds int) string {
	hours := durationInSeconds / 3600
	minutes := (durationInSeconds % 3600) / 60
	seconds := durationInSeconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
