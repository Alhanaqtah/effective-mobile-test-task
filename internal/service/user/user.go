package user

import (
	"errors"
	"log/slog"
	"time-tracker/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrExists       = errors.New("user already exists")
)

type Storage interface {
}

type PeopleRepo interface {
}

type Service struct {
	storage        Storage
	peopleInfoRepo PeopleRepo
	log            *slog.Logger
}

func New(storage Storage, peopleInfoRepo PeopleRepo, log *slog.Logger) *Service {
	return &Service{
		storage:        storage,
		peopleInfoRepo: peopleInfoRepo,
		log:            log,
	}
}

func (s *Service) GetUsers(page int, filter string) ([]models.User, error) {
	panic("unimplemented")
}

func (s *Service) RemoveUserByUUID(uuid string) error {
	panic("unimplemented")
}

func (s *Service) UpdateUserInfo(userInfo *models.User) (*models.User, error) {
	panic("unimplemented")
}

func (s *Service) CreateUser(passportSerie int, passportNumber int) error {
	panic("unimplemented")
}
