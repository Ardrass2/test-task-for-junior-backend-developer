package handlers

import (
	"time"

	periodictaskdomain "example.com/taskservice/internal/domain/periodictask"
)

type periodicTaskCreateDTO struct {
	Title            string                               `json:"title"`
	Description      string                               `json:"description"`
	RecurrenceType   periodictaskdomain.RecurrenceType     `json:"recurrence_type"`
	RecurrenceParams periodictaskdomain.RecurrenceParams   `json:"recurrence_params"`
	StartDate        string                               `json:"start_date"`
	EndDate          *string                              `json:"end_date,omitempty"`
}

type periodicTaskUpdateDTO struct {
	Title            string                               `json:"title"`
	Description      string                               `json:"description"`
	RecurrenceType   periodictaskdomain.RecurrenceType     `json:"recurrence_type"`
	RecurrenceParams periodictaskdomain.RecurrenceParams   `json:"recurrence_params"`
	StartDate        string                               `json:"start_date"`
	EndDate          *string                              `json:"end_date,omitempty"`
	IsActive         bool                                 `json:"is_active"`
}

type periodicTaskResponseDTO struct {
	ID               int64                                `json:"id"`
	Title            string                               `json:"title"`
	Description      string                               `json:"description"`
	RecurrenceType   periodictaskdomain.RecurrenceType     `json:"recurrence_type"`
	RecurrenceParams periodictaskdomain.RecurrenceParams   `json:"recurrence_params"`
	StartDate        string                               `json:"start_date"`
	EndDate          *string                              `json:"end_date,omitempty"`
	IsActive         bool                                 `json:"is_active"`
	CreatedAt        time.Time                            `json:"created_at"`
	UpdatedAt        time.Time                            `json:"updated_at"`
}

func newPeriodicTaskResponseDTO(pt *periodictaskdomain.PeriodicTask) periodicTaskResponseDTO {
	dto := periodicTaskResponseDTO{
		ID:               pt.ID,
		Title:            pt.Title,
		Description:      pt.Description,
		RecurrenceType:   pt.RecurrenceType,
		RecurrenceParams: pt.RecurrenceParams,
		StartDate:        pt.StartDate.Format("2006-01-02"),
		IsActive:         pt.IsActive,
		CreatedAt:        pt.CreatedAt,
		UpdatedAt:        pt.UpdatedAt,
	}
	if pt.EndDate != nil {
		formatted := pt.EndDate.Format("2006-01-02")
		dto.EndDate = &formatted
	}
	return dto
}
