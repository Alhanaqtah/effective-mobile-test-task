package user

import (
	"errors"
	"log/slog"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Storage interface {
}

type PeopleRepo interface {
}

// TODO: inject proplr repo
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
