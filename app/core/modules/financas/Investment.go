package financas

import (
	"time"
	"fmt"
)

type Investment struct {}
func (i Investment) GetDates(initialDate, finalDate string) ([]time.Time, error) {
	layout := "2006-01-02"
	start, err := time.Parse(layout, initialDate)
	if err != nil {
		return nil, fmt.Errorf("invalid initial date: %v", err)
	}
	end, err := time.Parse(layout, finalDate)
	if err != nil {
		return nil, fmt.Errorf("invalid final date: %v", err)
	}

	if end.Before(start) {
		return nil, fmt.Errorf("final date must be after initial date")
	}
    months := []time.Time{}
	for start.Before(end) {
		start = start.AddDate(0, 1, 0)
        months = append(months, start)
	}

	return months, nil
}
func (i Investment) MonthsBetweenDates(initialDate, finalDate string) (int, error) {
	layout := "2006-01-02"
	start, err := time.Parse(layout, initialDate)
	if err != nil {
		return 0, fmt.Errorf("invalid initial date: %v", err)
	}
	end, err := time.Parse(layout, finalDate)
	if err != nil {
		return 0, fmt.Errorf("invalid final date: %v", err)
	}

	if end.Before(start) {
		return 0, fmt.Errorf("final date must be after initial date")
	}

	months := 0
	for start.Before(end) {
		start = start.AddDate(0, 1, 0)
		months++
	}

	return months, nil
}
