package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"time-tracker/internal/config"
	"time-tracker/internal/models"
	"time-tracker/internal/repository"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(cfg *config.Storage) (*Storage, error) {
	const op = "repository.postgres.New"

	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	db := stdlib.OpenDB(*pool.Config().ConnConfig)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	m.Up()

	return &Storage{pool: pool}, nil
}

func (s *Storage) GetUsers(ctx context.Context, limit int, offset int, filter string) ([]models.User, error) {
	const op = "repository.postgresGetUsers"

	filter = "%" + filter + "%"

	rows, err := s.pool.Query(ctx, `
	SELECT id, name, surname, patronymic, address, passport_serie, passport_number 
	FROM users 
	WHERE name ILIKE $1 OR surname ILIKE $1 OR patronymic ILIKE $1 OR address ILIKE $1 OR CAST(passport_serie AS TEXT) ILIKE $1 OR CAST(passport_number AS TEXT) ILIKE $1 
	ORDER BY id 
	LIMIT $2 OFFSET $3
	`, filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Address, &user.PassportSerie, &user.PassportNumber)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (s *Storage) UpdateUser(ctx context.Context, fields []string, values []string) (*models.User, error) {
	const op = "repository.postgres.UpdateUser"

	f := strings.Join(fields[:len(fields)], ", ")

	v := make([]interface{}, len(values))
	for i, j := range values {
		v[i] = j
	}

	q := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d RETURNING id, name, surname, patronymic, address, passport_serie, passport_number", f, len(values))

	row := s.pool.QueryRow(ctx, q, v...)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Address, &user.PassportSerie, &user.PassportNumber)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) RemoveUser(ctx context.Context, uuid string) error {
	const op = "repository.postgres.RemoveUser"

	ct, err := s.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
	}

	return nil
}

func (s *Storage) Close() {
	s.pool.Close()
}
