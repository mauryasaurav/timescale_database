package repository

import (
	"context"
	"time"

	"github.com/mauryasaurav/timescale_database/models"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// UserRepository -
type UserRepository struct {
	conn *sqlx.DB
}

// NewUserRepository -
func NewUserRepository(conn *sqlx.DB) *UserRepository {
	return &UserRepository{conn}
}

// Save -
func (r *UserRepository) Save(ctx context.Context, user *models.UserCreateAndUpdate) error {
	query := "insert into users (time, req_bytes, username, pass, email) " +
		"values ($1, $2, $3, $4, $5);"
	_, err := r.conn.ExecContext(
		ctx,
		query,
		time.Now(),
		user.ReqBytes,
		user.UserName,
		user.Pass,
		user.Email,
	)

	if err != nil {
		return errors.Wrap(err, "error inserting user into database")
	}

	return nil
}

// GetUserBytesByFilter -
func (r *UserRepository) GetUserBytesByFilter(ctx context.Context, name string) (*models.UserResponse, error) {
	user := models.UserResponse{}
	query := "SELECT * FROM users_hourly where username = $1"
	row := r.conn.QueryRowxContext(ctx, query, name)

	if err := row.StructScan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
