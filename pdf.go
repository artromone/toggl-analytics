package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

func GenerateTablePdf(columns []string, rows [][]string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	colWidth, rowHeight := 45.0, 10.0

	for _, col := range columns {
		pdf.CellFormat(colWidth, rowHeight, col, "1", 0, "", false, 0, "")
	}
	pdf.Ln(rowHeight)

	for _, row := range rows {
		for _, col := range row {
			pdf.CellFormat(colWidth, rowHeight, col, "1", 0, "", false, 0, "")
		}
		pdf.Ln(rowHeight)
	}

	err := pdf.OutputFileAndClose("table.pdf")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("PDF created successfully")
}
