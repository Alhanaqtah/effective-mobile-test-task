package postgres

import (
	"context"
	"fmt"

	"time-tracker/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Storage) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	))
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &Storage{pool: pool}, nil
}

func (s *Storage) UpdateUser(ctx context.Context, fields, values []string) error {
	const op = "repository.postgres.UpdateUser"

	s.pool.QueryRow(ctx, fmt.Sprintf(`UPDATE users SET %s`), values)

	return nil
}

func (s *Storage) RemoveUser(ctx context.Context, uuid string) error {
	const op = "repository.postgres.RemoveUser"

	_, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Close() {
	s.pool.Close()
}
