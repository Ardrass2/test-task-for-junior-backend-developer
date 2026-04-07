package periodictask

import (
	"encoding/json"
	"time"
)

type RecurrenceType string

const (
	RecurrenceDaily         RecurrenceType = "daily"
	RecurrenceMonthly       RecurrenceType = "monthly"
	RecurrenceSpecificDates RecurrenceType = "specific_dates"
	RecurrenceEvenOdd       RecurrenceType = "even_odd"
)

func (r RecurrenceType) Valid() bool {
	switch r {
	case RecurrenceDaily, RecurrenceMonthly, RecurrenceSpecificDates, RecurrenceEvenOdd:
		return true
	default:
		return false
	}
}

type DayParity string

const (
	ParityEven DayParity = "even"
	ParityOdd  DayParity = "odd"
)

type RecurrenceParams struct {
	Interval *int       `json:"interval,omitempty"`
	Day      *int       `json:"day,omitempty"`
	Dates    []string   `json:"dates,omitempty"`
	Parity   *DayParity `json:"parity,omitempty"`
}

type PeriodicTask struct {
	ID               int64            `json:"id"`
	Title            string           `json:"title"`
	Description      string           `json:"description"`
	RecurrenceType   RecurrenceType   `json:"recurrence_type"`
	RecurrenceParams RecurrenceParams `json:"recurrence_params"`
	StartDate        time.Time        `json:"start_date"`
	EndDate          *time.Time       `json:"end_date,omitempty"`
	IsActive         bool             `json:"is_active"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
}

func (p *RecurrenceParams) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}

func ParseRecurrenceParams(data []byte) (RecurrenceParams, error) {
	var params RecurrenceParams
	if err := json.Unmarshal(data, &params); err != nil {
		return RecurrenceParams{}, err
	}
	return params, nil
}
