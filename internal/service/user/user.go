package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"time-tracker/internal/lib/logger/sl"
	"time-tracker/internal/models"
	"time-tracker/internal/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrExists       = errors.New("user already exists")
	ErrEmptyBody    = errors.New("request body is empty")
)

type Storage interface {
	RemoveUser(ctx context.Context, uuid string) error
	UpdateUser(ctx context.Context, fields []string, values []string) (*models.User, error)
	GetUsers(ctx context.Context, limit, offset int, filter string) ([]models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
}

type ExternalAPI interface {
	GetUserInfo(passportSerie, passportNumber int) (*models.User, error)
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

func (s *Service) CreateUser(ctx context.Context, passportSerie, passportNumber int) (*models.User, error) {
	const op = "service.user.CreateUser"

	log := s.log.With(slog.String("op", op))

	u, err := s.externalAPI.GetUserInfo(passportSerie, passportNumber)
	if err != nil {
		log.Error("failed to get user info from external api", sl.Error(err))
		return nil, err
	}

	u.PassportSerie = passportSerie
	u.PassportNumber = passportNumber

	user, err := s.storage.CreateUser(ctx, u)
	if err != nil {
		log.Error("failed to save user in storage", sl.Error(err))
		if errors.Is(err, repository.ErrExists) {
			return nil, ErrExists
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUsers(ctx context.Context, page int, filter string) ([]models.User, error) {
	const op = "service.user.GetUsers"

	log := s.log.With(slog.String("op", op))

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
	const op = "service.user.UpdateUserInfo"

	log := s.log.With(slog.String("op", op))

	var fields []string
	var values []string
	order := 1

	log.Debug("request body", slog.Any("body", userInfo))

	if userInfo.Name != "" {
		fields = append(fields, fmt.Sprintf("name = $%d", order))
		values = append(values, userInfo.Name)
		order++
	}
	if userInfo.Surname != "" {
		fields = append(fields, fmt.Sprintf("surname = $%d", order))
		values = append(values, userInfo.Surname)
		order++
	}
	if userInfo.Patronymic != "" {
		fields = append(fields, fmt.Sprintf("patronymic = $%d", order))
		values = append(values, userInfo.Patronymic)
		order++
	}
	if userInfo.Address != "" {
		fields = append(fields, fmt.Sprintf("address = $%d", order))
		values = append(values, userInfo.Address)
	}

	if len(values) == 0 {
		log.Debug("request body is empty")
		return nil, ErrEmptyBody
	}

	values = append(values, userInfo.ID)

	log.Debug("new user info", slog.Any("fields", fields), slog.Any("values", values))

	user, err := s.storage.UpdateUser(ctx, fields, values)
	if err != nil {
		log.Error("failed to update user info", sl.Error(err))
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) RemoveUserByUUID(ctx context.Context, uuid string) error {
	const op = "service.user.RemoveUserByUUID"

	log := s.log.With(slog.String("op", op))

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
