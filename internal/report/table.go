package report

import (
	"fmt"
)

type TableRow map[string]interface{}
type Table map[int]TableRow

func (t Table) AddRowById(id int, user string, duration int, sum float64, client, task string, taskTrackerId int) {
	t[id] = TableRow{
		"User":         user,
		"Duration":     duration,
		"Sum":          sum,
		"Client":       client,
		"Task":         task,
		"Vikunja link": taskTrackerId,
	}
}

func (t Table) AddRow(user string, duration int, sum float64, client, task string, taskTrackerId int) int {
	rows, _ := t.GetDimensions()
	t.AddRowById(rows+1, user, duration, sum, client, task, taskTrackerId)
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
	row, exists := t[id]
	if !exists {
		fmt.Printf("Row with ID %d not found\n", id)
		return nil
	}

	if _, ok := row[column]; !ok {
		fmt.Printf("Column %s not found]n", column)
		return nil
	}

	return t[id][column]
}

func (t Table) UpdateRow(id int, column string, value interface{}) error {
	row, exists := t[id]
	if !exists {
		return fmt.Errorf("Row with ID %d not found", id)
	}

	if _, ok := row[column]; !ok {
		return fmt.Errorf("Column %s not found", column)
	}

	row[column] = value
	return nil
}

func (t Table) DeleteRow(id int) {
	delete(t, id)
}
