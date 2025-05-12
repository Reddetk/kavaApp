// internal/infrastructure/postgres/user_repository.go
package postgres

import (
	"context"
	"database/sql"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	query := `INSERT INTO users (email, phone, age, gender, city) VALUES ($1, $2, $3, $4, $5) RETURNING id, registration_date, last_activity`
	return r.db.QueryRowContext(ctx, query,
		user.Email,
		user.Phone,
		user.Age,
		user.Gender,
		user.City,
	).Scan(&user.ID, &user.RegistrationDate, &user.LastActivity)
}

func (r *UserRepository) Get(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `SELECT id, email, phone, age, gender, city, registration_date, last_activity FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var user entities.User
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&user.Age,
		&user.Gender,
		&user.City,
		&user.RegistrationDate,
		&user.LastActivity,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	query := `UPDATE users SET email = $1, phone = $2, age = $3, gender = $4, city = $5, last_activity = NOW() WHERE id = $6`
	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Phone,
		user.Age,
		user.Gender,
		user.City,
		user.ID,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepository) List(ctx context.Context, limit int, offset int) ([]*entities.User, error) {
	query := `SELECT id, email, phone, age, gender, city, registration_date, last_activity FROM users ORDER BY registration_date DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Phone,
			&user.Age,
			&user.Gender,
			&user.City,
			&user.RegistrationDate,
			&user.LastActivity,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, rows.Err()
}
