package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"taskforge/models"
)

const userColumns = `id, username, email, role, password_hash`

func scanUser(row pgx.Row) (models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.PasswordHash)
	return u, err
}

func CreateUser(ctx context.Context, pool *pgxpool.Pool, user *models.User) error {
	if user.Role == "" {
		user.Role = "user"
	}
	err := pool.QueryRow(ctx,
		`INSERT INTO users (username, email, role, password_hash)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		user.Username, user.Email, user.Role, user.PasswordHash,
	).Scan(&user.ID)
	if isUniqueViolation(err) {
		return ErrAlreadyExists
	}
	return err
}

func GetUserByID(ctx context.Context, pool *pgxpool.Pool, id int) (models.User, error) {
	u, err := scanUser(pool.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = $1`, id))
	if errors.Is(err, pgx.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}

func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (models.User, error) {
	u, err := scanUser(pool.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE email = $1`, email))
	if errors.Is(err, pgx.ErrNoRows) {
		return u, ErrNotFound
	}
	return u, err
}

func ListUsers(ctx context.Context, pool *pgxpool.Pool) ([]models.User, error) {
	rows, err := pool.Query(ctx, `SELECT `+userColumns+` FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func UpdateUser(ctx context.Context, pool *pgxpool.Pool, user models.User) error {
	err := execOne(ctx, pool,
		`UPDATE users SET username = $1, email = $2, role = $3 WHERE id = $4`,
		user.Username, user.Email, user.Role, user.ID)
	if isUniqueViolation(err) {
		return ErrAlreadyExists
	}
	return err
}

func UpdateUserPassword(ctx context.Context, pool *pgxpool.Pool, id int, passwordHash string) error {
	return execOne(ctx, pool,
		`UPDATE users SET password_hash = $1 WHERE id = $2`,
		passwordHash, id)
}

func DeleteUser(ctx context.Context, pool *pgxpool.Pool, id int) error {
	return execOne(ctx, pool, `DELETE FROM users WHERE id = $1`, id)
}
