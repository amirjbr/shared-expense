package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/amirjbr/shared-expense/internal/account/core/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) (*UserRepo, error) {

	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &UserRepo{
		DB: db,
	}, nil
}

// ---------------------------------------------------------
// Create User
// ---------------------------------------------------------

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) error {

	const query = `
		INSERT INTO users (
			id,
			first_name,
			last_name,
			username,
			password,
			email,
			phone_number,
			created_at,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9
		)
	`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Password,
		user.Email,
		user.PhoneNumber,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf(
			"insert user: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get create user rows affected: %w",
			err,
		)
	}

	if rowsAffected != 1 {
		return errors.New(
			"user was not created",
		)
	}

	return nil
}

// ---------------------------------------------------------
// Get User By ID
// ---------------------------------------------------------

//no usage for now delete if not necessary

func (r *UserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {

	const query = `
		SELECT
			id,
			first_name,
			last_name,
			username,
			password,
			email,
			phone_number,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var user entity.User

	err := r.DB.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"get user by ID: %w",
			err,
		)
	}

	return user, nil
}

// ---------------------------------------------------------
// Get User By Username
// ---------------------------------------------------------

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {

	const query = `
		SELECT
			id,
			first_name,
			last_name,
			username,
			password,
			email,
			phone_number,
			created_at,
			updated_at
		FROM users
		WHERE username = $1
	`

	var user entity.User

	err := r.DB.QueryRowContext(
		ctx,
		query,
		username,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return entity.User{}, fmt.Errorf(
			"get user by username: %w",
			err,
		)
	}

	return user, nil
}

// ---------------------------------------------------------
// Update User
// ---------------------------------------------------------

func (r *UserRepo) UpdateUser(ctx context.Context, user entity.User) error {

	const query = `
		UPDATE users
		SET
			first_name = $1,
			last_name = $2,
			username = $3,
			password = $4,
			email = $5,
			phone_number = $6,
			updated_at = $7
		WHERE id = $8
	`

	result, err := r.DB.ExecContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Password,
		user.Email,
		user.PhoneNumber,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf(
			"update user: %w",
			err,
		)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(
			"get update user rows affected: %w",
			err,
		)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ---------------------------------------------------------
// Get Multiple Users
// ---------------------------------------------------------

func (r *UserRepo) GetUsersByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.User, error) {

	if len(ids) == 0 {
		return []entity.User{}, nil
	}

	const query = `
		SELECT
			id,
			first_name,
			last_name,
			username,
			password,
			email,
			phone_number,
			created_at,
			updated_at
		FROM users
		WHERE id = ANY($1::uuid[])
	`

	/*
		Conversion is kept inside the Postgres adapter because
		lib/pq handles []string arrays conveniently.
	*/
	idStrings := make(
		[]string,
		len(ids),
	)

	for i, id := range ids {
		idStrings[i] = id.String()
	}

	rows, err := r.DB.QueryContext(
		ctx,
		query,
		pq.Array(idStrings),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get users by IDs: %w",
			err,
		)
	}

	defer rows.Close()

	users := make(
		[]entity.User,
		0,
		len(ids),
	)

	for rows.Next() {
		var user entity.User

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.PhoneNumber,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"scan user: %w",
				err,
			)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate users: %w",
			err,
		)
	}

	return users, nil
}
