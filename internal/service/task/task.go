package task

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"time-tracker/internal/lib/logger/sl"
	"time-tracker/internal/models"

	"github.com/google/uuid"
)

var (
	ErrInvalidDateRange = errors.New("invalid date range")
	ErrInvalidUUID      = errors.New("invalid uuid format")
	ErrInvalidDate      = errors.New("invalid date format")
)

type Storage interface {
	GetTasksInRange(ctx context.Context, userUUID string, startDate, endDate time.Time) ([]models.Task, error)
}

type Service struct {
	storage Storage
	log     *slog.Logger
}

func New(storage Storage, log *slog.Logger) *Service {
	return &Service{
		storage: storage,
		log:     log,
	}
}

func (s *Service) GetTasksInRange(ctx context.Context, userUUID, startDate, endDate string) ([]models.Task, error) {
	const op = "service.task.GetTasksInRange"

	log := s.log.With(slog.String("op", op))

	log.Debug("validating input parameters", slog.String("userUUID", userUUID), slog.String("startDate", startDate), slog.String("endDate", endDate))

	// Validate userUUID
	_, err := uuid.Parse(userUUID)
	if err != nil {
		log.Error("invalid userUUID", sl.Error(err))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidUUID)
	}

	// Parse startDate and endDate
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		log.Error("invalid startDate", sl.Error(err))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidDate)
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		log.Error("invalid endDate", sl.Error(err))
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidDate)
	}

	if start.After(end) {
		log.Error("startDate is after endDate")
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidDateRange)
	}

	log.Debug("fetching tasks from storage", slog.String("userUUID", userUUID), slog.Time("startDate", start), slog.Time("endDate", end))

	tasks, err := s.storage.GetTasksInRange(ctx, userUUID, start, end)
	if err != nil {
		log.Error("failed to fetch tasks from storage", sl.Error(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("tasks fetched successfully", slog.Int("tasksCount", len(tasks)))

	return tasks, nil
}
