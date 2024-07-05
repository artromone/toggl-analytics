package main

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
)

func GenerateTablePdf(columns []string, rows [][]string) {
	pdf := gofpdf.New("P", "mm", "A4", "font/")
	// fontPath := "fonts/Montserrat-Bold.ttf"
	// pdf.AddFont("Montserrat-Bold", "", fontPath)

	pdf.AddFont("Helvetica", "", "helvetica_1251.json")
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 12)

	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")

	// pdf.Cell(15, 50, tr("русский текст"))

	colWidth, rowHeight := 45.0, 10.0

	for _, col := range columns {
		pdf.CellFormat(colWidth, rowHeight, tr(col), "1", 0, "", false, 0, "")
		// println(col)
	}
	pdf.Ln(rowHeight)

	for _, row := range rows {
		for _, col := range row {
			pdf.CellFormat(colWidth, rowHeight, tr(col), "1", 0, "", false, 0, "")
			// println(col)
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
