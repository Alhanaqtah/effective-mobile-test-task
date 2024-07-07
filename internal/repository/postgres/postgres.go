package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"

	"time-tracker/internal/config"
	"time-tracker/internal/models"

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

func (s *Storage) GetUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	const op = "repository.postgres.GetUsers"

	rows, err := s.pool.Query(ctx, "SELECT id, name, surname, patronymic, address, passport_serie, passport_number FROM users ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Address, &user.PassportSerie, &user.PassportNumber)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) UpdateUser(ctx context.Context, fields, values []string) (*models.User, error) {
	const op = "repository.postgres.UpdateUser"

	q := fmt.Sprintf("UPDATE users SET %s RETURNING id, name, surname, patronymic, address, password_serie, passwort_number", strings.Join(fields, ", "))

	row := s.pool.QueryRow(ctx, q, values)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Address, &user.PassportSerie, &user.PassportNumber)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
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
