package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	periodictaskdomain "example.com/taskservice/internal/domain/periodictask"
)

type PeriodicTaskRepository struct {
	pool *pgxpool.Pool
}

func NewPeriodicTaskRepository(pool *pgxpool.Pool) *PeriodicTaskRepository {
	return &PeriodicTaskRepository{pool: pool}
}

func (r *PeriodicTaskRepository) Create(ctx context.Context, pt *periodictaskdomain.PeriodicTask) (*periodictaskdomain.PeriodicTask, error) {
	paramsJSON, err := json.Marshal(pt.RecurrenceParams)
	if err != nil {
		return nil, err
	}

	const query = `
		INSERT INTO periodic_tasks (title, description, recurrence_type, recurrence_params, start_date, end_date, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, title, description, recurrence_type, recurrence_params, start_date, end_date, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		pt.Title, pt.Description, pt.RecurrenceType, paramsJSON,
		pt.StartDate, pt.EndDate, pt.IsActive, pt.CreatedAt, pt.UpdatedAt,
	)

	return scanPeriodicTask(row)
}

func (r *PeriodicTaskRepository) GetByID(ctx context.Context, id int64) (*periodictaskdomain.PeriodicTask, error) {
	const query = `
		SELECT id, title, description, recurrence_type, recurrence_params, start_date, end_date, is_active, created_at, updated_at
		FROM periodic_tasks
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	pt, err := scanPeriodicTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, periodictaskdomain.ErrNotFound
		}
		return nil, err
	}
	return pt, nil
}

func (r *PeriodicTaskRepository) Update(ctx context.Context, pt *periodictaskdomain.PeriodicTask) (*periodictaskdomain.PeriodicTask, error) {
	paramsJSON, err := json.Marshal(pt.RecurrenceParams)
	if err != nil {
		return nil, err
	}

	const query = `
		UPDATE periodic_tasks
		SET title = $1, description = $2, recurrence_type = $3, recurrence_params = $4,
			start_date = $5, end_date = $6, is_active = $7, updated_at = $8
		WHERE id = $9
		RETURNING id, title, description, recurrence_type, recurrence_params, start_date, end_date, is_active, created_at, updated_at
	`

	row := r.pool.QueryRow(ctx, query,
		pt.Title, pt.Description, pt.RecurrenceType, paramsJSON,
		pt.StartDate, pt.EndDate, pt.IsActive, pt.UpdatedAt, pt.ID,
	)

	updated, err := scanPeriodicTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, periodictaskdomain.ErrNotFound
		}
		return nil, err
	}
	return updated, nil
}

func (r *PeriodicTaskRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM periodic_tasks WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return periodictaskdomain.ErrNotFound
	}
	return nil
}

func (r *PeriodicTaskRepository) List(ctx context.Context) ([]periodictaskdomain.PeriodicTask, error) {
	const query = `
		SELECT id, title, description, recurrence_type, recurrence_params, start_date, end_date, is_active, created_at, updated_at
		FROM periodic_tasks
		ORDER BY id DESC
	`
	return r.queryPeriodicTasks(ctx, query)
}

func (r *PeriodicTaskRepository) ListActive(ctx context.Context) ([]periodictaskdomain.PeriodicTask, error) {
	const query = `
		SELECT id, title, description, recurrence_type, recurrence_params, start_date, end_date, is_active, created_at, updated_at
		FROM periodic_tasks
		WHERE is_active = TRUE
		ORDER BY id
	`
	return r.queryPeriodicTasks(ctx, query)
}

func (r *PeriodicTaskRepository) queryPeriodicTasks(ctx context.Context, query string) ([]periodictaskdomain.PeriodicTask, error) {
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]periodictaskdomain.PeriodicTask, 0)
	for rows.Next() {
		pt, err := scanPeriodicTask(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, *pt)
	}
	return result, rows.Err()
}

type periodicTaskScanner interface {
	Scan(dest ...any) error
}

func scanPeriodicTask(scanner periodicTaskScanner) (*periodictaskdomain.PeriodicTask, error) {
	var (
		pt             periodictaskdomain.PeriodicTask
		recurrenceType string
		paramsJSON     []byte
		startDate      time.Time
		endDate        *time.Time
	)

	if err := scanner.Scan(
		&pt.ID, &pt.Title, &pt.Description,
		&recurrenceType, &paramsJSON,
		&startDate, &endDate,
		&pt.IsActive, &pt.CreatedAt, &pt.UpdatedAt,
	); err != nil {
		return nil, err
	}

	pt.RecurrenceType = periodictaskdomain.RecurrenceType(recurrenceType)
	pt.StartDate = startDate
	pt.EndDate = endDate

	if err := json.Unmarshal(paramsJSON, &pt.RecurrenceParams); err != nil {
		return nil, err
	}

	return &pt, nil
}
