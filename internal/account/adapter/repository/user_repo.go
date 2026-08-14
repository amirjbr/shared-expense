package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/amirjbr/shared-expense/internal/account/core/entity"
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

func (r *UserRepo) CreateUser(ctx context.Context, user entity.User) (string, error) {
	query := `INSERT INTO users(id,first_name,last_name,username,password,
                  				email,phone_number,created_at,updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8,$9);`

	result, err := r.DB.ExecContext(ctx, query, user.ID, user.FirstName, user.LastName, user.Username,
		user.Password, user.Email, user.PhoneNumber, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return "", err
	}
	rowAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowAffected != 1 {
		return "", errors.New("user not created")
	}

	return user.ID.String(), nil

}
func (r *UserRepo) GetUserByID(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	query := `Select * from users where id = $1;`

	err := r.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username,
		&user.Password, &user.Email, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	query := `Select * from users where username = $1;`

	err := r.DB.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Username,
		&user.Password, &user.Email, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
func (r *UserRepo) UpdateUser(ctx context.Context, user entity.User) error {
	// TODO update need to fix
	query := `UPDATE users set first_name = $1, last_name = $2,username = $3 , password = $4,email = $5,
                 phonenumber = $6, updated_at = $7 where id = $8;`
	fmt.Println(query)
	return nil
}
