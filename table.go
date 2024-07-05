package main

import (
	"fmt"
)

type TableRow map[string]interface{}
type Table map[int]TableRow

func (t Table) AddRowById(id int, user string, duration int, sum float64, client, task string) {
	t[id] = TableRow{
		"User":     user,
		"Duration": duration,
		"Sum":      sum,
		"Client":   client,
		"Task":     task,
	}
}

func (t Table) AddRow(user string, duration int, sum float64, client, task string) int {
	rows, _ := t.GetDimensions()
	t.AddRowById(rows+1, user, duration, sum, client, task)
	return rows + 1
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

func (t Table) Get(id int, column string) interface{} {
	return t[id][column]
} // TODO

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

// func (t Table) PrintTable() {
// 	fmt.Printf("%-5s %-10s %-20s %-10s %-15s %-15s\n", "ID", "User", "Duration", "Sum", "Client", "Task")
// 	for id, row := range t {
// 		fmt.Printf("%-5d %-10s %-20s %-10.2f %-15s %-15s\n",
// 			id,
// 			row["User"].(string),
// 			row["Duration"].(int),
// 			row["Sum"].(float64),
// 			row["Client"].(string),
// 			row["Task"].(string))
// 	}
// }
