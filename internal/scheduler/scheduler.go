package scheduler

import (
	"context"
	"log/slog"
	"time"

	periodictaskusecase "example.com/taskservice/internal/usecase/periodictask"
)

type Scheduler struct {
	usecase  periodictaskusecase.Usecase
	logger   *slog.Logger
	interval time.Duration
}

func New(usecase periodictaskusecase.Usecase, logger *slog.Logger, interval time.Duration) *Scheduler {
	return &Scheduler{
		usecase:  usecase,
		logger:   logger,
		interval: interval,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	s.logger.Info("scheduler started", "interval", s.interval)
	s.generate(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopped")
			return
		case <-ticker.C:
			s.generate(ctx)
		}
	}
}

func (s *Scheduler) generate(ctx context.Context) {
	today := time.Now().UTC()
	created, err := s.usecase.GenerateTasksForDate(ctx, today)
	if err != nil {
		s.logger.Error("scheduler: failed to generate tasks", "error", err)
		return
	}
	if created > 0 {
		s.logger.Info("scheduler: generated tasks", "count", created, "date", today.Format("2006-01-02"))
	}
}
