package repository

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrExists       = errors.New("user already exists")
)
