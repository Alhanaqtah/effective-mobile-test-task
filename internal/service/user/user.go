package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"time-tracker/internal/lib/logger/sl"
	"time-tracker/internal/models"
	"time-tracker/internal/repository"

	"github.com/go-chi/chi/v5/middleware"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrExists       = errors.New("user already exists")
)

type Storage interface {
	RemoveUser(ctx context.Context, uuid string) error
	UpdateUser(ctx context.Context, fields, values []string) (*models.User, error)
	GetUsers(ctx context.Context, limit, offset int, filter string) ([]models.User, error)
}

type ExternalAPI interface {
}

type Service struct {
	storage     Storage
	externalAPI ExternalAPI
	log         *slog.Logger
}

func New(storage Storage, externalAPI ExternalAPI, log *slog.Logger) *Service {
	return &Service{
		storage:     storage,
		externalAPI: externalAPI,
		log:         log,
	}
}

func (s *Service) CreateUser(ctx context.Context, passportSerie int, passportNumber int) error {
	panic("")
}

func (s *Service) GetUsers(ctx context.Context, page int, filter string) ([]models.User, error) {
	const op = "service.user.GetUsers"

	log := s.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(ctx)),
	)

	const limit = 10
	offset := (page - 1) * limit

	users, err := s.storage.GetUsers(ctx, limit, offset, filter)
	if err != nil {
		log.Error("error while getting users", sl.Error(err))
		return nil, err
	}

	return users, nil
}

func (s *Service) UpdateUserInfo(ctx context.Context, userInfo *models.User) (*models.User, error) {
	const op = "service.user.RemoveUserByUUID"

	log := s.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(ctx)),
	)

	var fields []string
	var values []string
	order := 1

	if userInfo.Name != "" {
		fields = append(fields, fmt.Sprintf("name = $%d", order))
		values = append(values, userInfo.Name)
		order++
	} else if userInfo.Surname != "" {
		fields = append(fields, fmt.Sprintf("surname = $%d", order))
		values = append(values, userInfo.Surname)
		order++
	} else if userInfo.Patronymic != "" {
		fields = append(fields, fmt.Sprintf("patronymic = $%d", order))
		values = append(values, userInfo.Patronymic)
	}

	log.Debug("new user info", slog.Any("fields", fields), slog.Any("values", values))

	user, err := s.storage.UpdateUser(ctx, fields, values)
	if err != nil {
		log.Error("failed to update user info", sl.Error(err))
		return nil, err
	}

	return user, nil
}

func (s *Service) RemoveUserByUUID(ctx context.Context, uuid string) error {
	const op = "service.user.RemoveUserByUUID"

	log := s.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(ctx)),
	)

	err := s.storage.RemoveUser(ctx, uuid)
	if err != nil {
		log.Error("failed to remove user by uuid", sl.Error(err))
		if errors.Is(err, repository.ErrUserNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return nil
}
