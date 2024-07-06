package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"os"
	"sort"
	"strconv"
)

func GenerateTablePdf(columns []string, rows [][]string, colWidths map[int]float64) {
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
				value = tr(fmt.Sprintf("Перейти (id:%d)", i))
			}

			pdf.CellFormat(width, rowHeight, value, "1", 0, "", false, 0, linkStr)
		}
		pdf.Ln(rowHeight)
	}

	err := pdf.OutputFileAndClose("reports/table.pdf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("PDF created successfully")
}
