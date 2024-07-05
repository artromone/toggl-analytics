package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"sort"
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

	for _, row := range rows {
		for i, col := range row {
			width, exists := colWidths[i]
			if !exists {
				width = 20.0
			}
			pdf.CellFormat(width, rowHeight, tr(col), "1", 0, "", false, 0, "")
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
