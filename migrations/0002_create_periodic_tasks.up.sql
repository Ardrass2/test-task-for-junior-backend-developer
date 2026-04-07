CREATE TABLE IF NOT EXISTS periodic_tasks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    recurrence_type TEXT NOT NULL,
    recurrence_params JSONB NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS periodic_task_id BIGINT REFERENCES periodic_tasks(id) ON DELETE SET NULL;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS scheduled_date DATE;

CREATE UNIQUE INDEX IF NOT EXISTS idx_tasks_periodic_date
    ON tasks (periodic_task_id, scheduled_date)
    WHERE periodic_task_id IS NOT NULL;
