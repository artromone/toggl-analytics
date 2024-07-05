package main

import (
	"fmt"
	"time"
)

type TableRow map[string]interface{}
type Table map[int]TableRow

func (t Table) AddRowById(id int, user string, time time.Time, sum float64, project string) {
	t[id] = TableRow{
		"User":    user,
		"Time":    time,
		"Sum":     sum,
		"Project": project,
	}
}

func (t Table) AddRow(user string, time time.Time, sum float64, project string) {
	rows, _ := t.GetDimensions()
	t.AddRowById(rows+1, user, time, sum, project)
}

func (t Table) GetDimensions() (int, int) {
	rows := len(t)
	var columns int

	if rows > 0 {
		for _, row := range t {
			columns = len(row)
			break
		}
	}

	return rows, columns
}

func (t Table) UpdateRow(id int, column string, value interface{}) error {
	row, exists := t[id]
	if !exists {
		return fmt.Errorf("row with ID %d not found", id)
	}

	if _, ok := row[column]; !ok {
		return fmt.Errorf("column %s not found", column)
	}

	row[column] = value
	return nil
}

func (t Table) DeleteRow(id int) {
	delete(t, id)
}

func (t Table) PrintTable() {
	fmt.Printf("%-5s %-10s %-20s %-10s %-15s\n", "ID", "User", "Time", "Sum", "Project")
	for id, row := range t {
		fmt.Printf("%-5d %-10s %-20s %-10.2f %-15s\n",
			id,
			row["User"].(string),
			row["Time"].(time.Time).Format(time.RFC3339),
			row["Sum"].(float64),
			row["Project"].(string))
	}
}
