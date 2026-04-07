package periodictask

import (
	"context"
	"time"

	periodictaskdomain "example.com/taskservice/internal/domain/periodictask"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository interface {
	Create(ctx context.Context, pt *periodictaskdomain.PeriodicTask) (*periodictaskdomain.PeriodicTask, error)
	GetByID(ctx context.Context, id int64) (*periodictaskdomain.PeriodicTask, error)
	Update(ctx context.Context, pt *periodictaskdomain.PeriodicTask) (*periodictaskdomain.PeriodicTask, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]periodictaskdomain.PeriodicTask, error)
	ListActive(ctx context.Context) ([]periodictaskdomain.PeriodicTask, error)
}

type TaskRepository interface {
	CreateFromSchedule(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
}

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*periodictaskdomain.PeriodicTask, error)
	GetByID(ctx context.Context, id int64) (*periodictaskdomain.PeriodicTask, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*periodictaskdomain.PeriodicTask, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]periodictaskdomain.PeriodicTask, error)
	GenerateTasksForDate(ctx context.Context, date time.Time) (int, error)
}

type CreateInput struct {
	Title            string
	Description      string
	RecurrenceType   periodictaskdomain.RecurrenceType
	RecurrenceParams periodictaskdomain.RecurrenceParams
	StartDate        time.Time
	EndDate          *time.Time
}

type UpdateInput struct {
	Title            string
	Description      string
	RecurrenceType   periodictaskdomain.RecurrenceType
	RecurrenceParams periodictaskdomain.RecurrenceParams
	StartDate        time.Time
	EndDate          *time.Time
	IsActive         bool
}
