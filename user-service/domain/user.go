package domain

import (
	"context"
	"github.com/google/uuid"
	"time"
	"user-service/dto/authDto"
	"user-service/dto/userDto"
)

type User struct {
	ID            int64      `db:"id"`
	Uuid          uuid.UUID  `db:"uuid"`
	FullName      string     `db:"full_name"`
	Phone         string     `db:"phone"`
	Email         string     `db:"email"`
	Username      string     `db:"username"`
	Password      string     `db:"password"`
	EmailVerifyAt *time.Time `db:"email_verify_at"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}

type UserRepository interface {
	FindByUuid(ctx context.Context, uuid uuid.UUID) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	CheckEmailExistence(ctx context.Context, email string) error
	Insert(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type UserService interface {
	Register(ctx context.Context, req authDto.RegisterRequest) (authDto.RegisterResponse, error)
	Authenticate(ctx context.Context, req authDto.LoginRequest) (authDto.LoginResponse, error)
	Detail(ctx context.Context, uuid uuid.UUID) (userDto.UserResponse, error)
	Update(ctx context.Context, req userDto.UserUpdateRequest) (userDto.UserUpdateResponse, error)

	//ValidateToken()
	//ValidateOTP()
}
