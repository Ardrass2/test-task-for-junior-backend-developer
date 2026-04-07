package periodictask

import (
	"context"
	"fmt"
	"strings"
	"time"

	periodictaskdomain "example.com/taskservice/internal/domain/periodictask"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo     Repository
	taskRepo TaskRepository
	now      func() time.Time
}

func NewService(repo Repository, taskRepo TaskRepository) *Service {
	return &Service{
		repo:     repo,
		taskRepo: taskRepo,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*periodictaskdomain.PeriodicTask, error) {
	if err := validateCreateInput(input); err != nil {
		return nil, err
	}
	startDate, endDate, err := normalizeAndValidateDateRange(input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	now := s.now()
	pt := &periodictaskdomain.PeriodicTask{
		Title:            strings.TrimSpace(input.Title),
		Description:      strings.TrimSpace(input.Description),
		RecurrenceType:   input.RecurrenceType,
		RecurrenceParams: input.RecurrenceParams,
		StartDate:        startDate,
		EndDate:          endDate,
		IsActive:         true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return s.repo.Create(ctx, pt)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*periodictaskdomain.PeriodicTask, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*periodictaskdomain.PeriodicTask, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	if err := validateUpdateInput(input); err != nil {
		return nil, err
	}
	startDate, endDate, err := normalizeAndValidateDateRange(input.StartDate, input.EndDate)
	if err != nil {
		return nil, err
	}

	pt := &periodictaskdomain.PeriodicTask{
		ID:               id,
		Title:            strings.TrimSpace(input.Title),
		Description:      strings.TrimSpace(input.Description),
		RecurrenceType:   input.RecurrenceType,
		RecurrenceParams: input.RecurrenceParams,
		StartDate:        startDate,
		EndDate:          endDate,
		IsActive:         input.IsActive,
		UpdatedAt:        s.now(),
	}

	return s.repo.Update(ctx, pt)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]periodictaskdomain.PeriodicTask, error) {
	return s.repo.List(ctx)
}

func (s *Service) GenerateTasksForDate(ctx context.Context, date time.Time) (int, error) {
	active, err := s.repo.ListActive(ctx)
	if err != nil {
		return 0, fmt.Errorf("list active periodic tasks: %w", err)
	}

	created := 0
	for i := range active {
		pt := &active[i]
		if !pt.MatchesDate(date) {
			continue
		}

		scheduledDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		now := s.now()
		task := &taskdomain.Task{
			Title:          pt.Title,
			Description:    pt.Description,
			Status:         taskdomain.StatusNew,
			PeriodicTaskID: &pt.ID,
			ScheduledDate:  &scheduledDate,
			CreatedAt:      now,
			UpdatedAt:      now,
		}

		result, err := s.taskRepo.CreateFromSchedule(ctx, task)
		if err != nil {
			return created, fmt.Errorf("create task for periodic_task %d: %w", pt.ID, err)
		}
		if result != nil {
			created++
		}
	}

	return created, nil
}

func validateCreateInput(input CreateInput) error {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if !input.RecurrenceType.Valid() {
		return fmt.Errorf("%w: invalid recurrence_type", ErrInvalidInput)
	}
	return validateRecurrenceParams(input.RecurrenceType, input.RecurrenceParams)
}

func validateUpdateInput(input UpdateInput) error {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return fmt.Errorf("%w: title is required", ErrInvalidInput)
	}
	if !input.RecurrenceType.Valid() {
		return fmt.Errorf("%w: invalid recurrence_type", ErrInvalidInput)
	}
	return validateRecurrenceParams(input.RecurrenceType, input.RecurrenceParams)
}

func validateRecurrenceParams(rt periodictaskdomain.RecurrenceType, params periodictaskdomain.RecurrenceParams) error {
	switch rt {
	case periodictaskdomain.RecurrenceDaily:
		if params.Interval == nil || *params.Interval < 1 {
			return fmt.Errorf("%w: daily recurrence requires interval >= 1", ErrInvalidInput)
		}
	case periodictaskdomain.RecurrenceMonthly:
		if params.Day == nil || *params.Day < 1 || *params.Day > 30 {
			return fmt.Errorf("%w: monthly recurrence requires day between 1 and 30", ErrInvalidInput)
		}
	case periodictaskdomain.RecurrenceSpecificDates:
		if len(params.Dates) == 0 {
			return fmt.Errorf("%w: specific_dates recurrence requires at least one date", ErrInvalidInput)
		}
		for _, d := range params.Dates {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				return fmt.Errorf("%w: invalid date format %q, expected YYYY-MM-DD", ErrInvalidInput, d)
			}
		}
	case periodictaskdomain.RecurrenceEvenOdd:
		if params.Parity == nil {
			return fmt.Errorf("%w: even_odd recurrence requires parity (even or odd)", ErrInvalidInput)
		}
		if *params.Parity != periodictaskdomain.ParityEven && *params.Parity != periodictaskdomain.ParityOdd {
			return fmt.Errorf("%w: parity must be 'even' or 'odd'", ErrInvalidInput)
		}
	}
	return nil
}

func normalizeAndValidateDateRange(startDate time.Time, endDate *time.Time) (time.Time, *time.Time, error) {
	if startDate.IsZero() {
		return time.Time{}, nil, fmt.Errorf("%w: start_date is required", ErrInvalidInput)
	}

	normalizedStart := toUTCDate(startDate)
	var normalizedEnd *time.Time
	if endDate != nil {
		value := toUTCDate(*endDate)
		if value.Before(normalizedStart) {
			return time.Time{}, nil, fmt.Errorf("%w: end_date must be greater than or equal to start_date", ErrInvalidInput)
		}
		normalizedEnd = &value
	}

	return normalizedStart, normalizedEnd, nil
}

func toUTCDate(value time.Time) time.Time {
	return time.Date(value.UTC().Year(), value.UTC().Month(), value.UTC().Day(), 0, 0, 0, 0, time.UTC)
}
