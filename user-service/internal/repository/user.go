package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"user-service/domain"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUser(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{db: pool}
}

func (u userRepository) Insert(ctx context.Context, user *domain.User) error {
	_, err := u.db.Exec(ctx, `
		INSERT INTO users
		    (uuid, full_name, phone, email, username, password, email_verify_at, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.Uuid,
		user.FullName,
		user.Phone,
		user.Email,
		user.Username,
		user.Password,
		user.EmailVerifyAt,
		user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (u userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	row := u.db.QueryRow(ctx, `
		SELECT uuid, full_name, phone, email, username, password
		FROM users
		WHERE email = $1`, email)

	err := row.Scan(&user.Uuid, &user.FullName, &user.Phone, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound // Assuming domain.ErrNotFound is defined for this case
		}
		fmt.Println(&user)
		return domain.User{}, fmt.Errorf("error scanning user: %w", err)
	}
	return user, nil
}

func (u userRepository) CheckEmailExistence(ctx context.Context, email string) error {
	var exists bool
	err := u.db.QueryRow(ctx, `
        SELECT EXISTS (
            SELECT 1
            FROM users
            WHERE email = $1
        )`, email).Scan(&exists)

	if err != nil {
		return err // Handle database errors
	}

	if !exists {
		// Email doesn't exist, return nil
		return nil
	}

	// Email exists, return an error
	return errors.New("email already exists")
}

func (u userRepository) FindByUuid(ctx context.Context, uuid uuid.UUID) (domain.User, error) {
	var user domain.User

	row := u.db.QueryRow(ctx, `
		SELECT uuid, full_name, phone, email, username, password,email_verify_at, created_at
		FROM users
		WHERE uuid = $1`, uuid)

	err := row.Scan(&user.Uuid, &user.FullName, &user.Phone, &user.Email, &user.Username, &user.Password, &user.EmailVerifyAt, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound // Assuming domain.ErrNotFound is defined for this case
		}
		return domain.User{}, fmt.Errorf("error scanning user: %w", err)
	}

	return user, nil
}

func (u userRepository) Update(ctx context.Context, user *domain.User) error {
	_, err := u.db.Exec(ctx, `
		UPDATE users
		SET
			full_name = $2,
			phone = $3,
			email = $4,
			username = $5,
			password = $6,
			updated_at = $7
		WHERE
			uuid = $1
		`, user.Uuid, user.FullName, user.Phone, user.Email, user.Username, user.Password, user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}
