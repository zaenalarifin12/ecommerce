package service

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"log"
	"time"
	"user-service/domain"
	"user-service/dto/authDto"
	"user-service/dto/userDto"
	"user-service/internal/utils"
)

var jwtSecret = []byte("your_secret_key")

type userService struct {
	userRepository domain.UserRepository
}

func NewUser(userRepository domain.UserRepository) domain.UserService {
	return &userService{userRepository: userRepository}
}

func (u userService) Register(ctx context.Context, req authDto.RegisterRequest) (authDto.RegisterResponse, error) {

	err := u.userRepository.CheckEmailExistence(ctx, req.Email)

	if err != nil {
		return authDto.RegisterResponse{}, err
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		panic("Error hash password")
	}

	user := domain.User{
		Uuid:      uuid.New(),
		FullName:  req.FullName,
		Phone:     req.Phone,
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashPassword,
		CreatedAt: time.Now(),
	}

	err = u.userRepository.Insert(ctx, &user)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return authDto.RegisterResponse{}, err
	}

	return authDto.RegisterResponse{
		Uuid:     user.Uuid,
		FullName: user.FullName,
		Phone:    user.Phone,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (u userService) Authenticate(ctx context.Context, req authDto.LoginRequest) (authDto.LoginResponse, error) {
	ctx, parentSpan := tracer.Start(ctx, "User.AuthenticateV1")
	defer parentSpan.End()

	// check if user exist
	parentSpan.AddEvent("find email")
	user, err := u.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		parentSpan.RecordError(err)
		return authDto.LoginResponse{}, err
	}

	// compare password
	err = utils.VerifyPassword(user.Password, req.Password)

	if err != nil {
		parentSpan.RecordError(err)
		return authDto.LoginResponse{}, err
	}

	// Generate jwt token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = user.Uuid
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // token expired

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("Error signing token: %s", err.Error())
		parentSpan.RecordError(err)
		return authDto.LoginResponse{}, err
	}

	// return
	return authDto.LoginResponse{
		Token:    tokenString,
		Uuid:     user.Uuid,
		FullName: user.FullName,
		Phone:    user.Phone,
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (u userService) Detail(ctx context.Context, uuid uuid.UUID) (userDto.UserResponse, error) {

	//	find users by uuid
	user, err := u.userRepository.FindByUuid(ctx, uuid)

	if err != nil {
		return userDto.UserResponse{}, err
	}

	if user == (domain.User{}) {
		return userDto.UserResponse{}, domain.ErrNotFound
	}

	return userDto.UserResponse{
		Uuid:          user.Uuid,
		FullName:      user.FullName,
		Phone:         user.Phone,
		Email:         user.Email,
		Username:      user.Username,
		EmailVerifyAt: user.EmailVerifyAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}, nil

}

func (u userService) Update(ctx context.Context, req userDto.UserUpdateRequest) (userDto.UserUpdateResponse, error) {
	//	find users by uuid
	user, err := u.userRepository.FindByUuid(ctx, req.Uuid)

	if err != nil {
		return userDto.UserUpdateResponse{}, err
	}

	if user == (domain.User{}) {
		return userDto.UserUpdateResponse{}, domain.ErrNotFound
	}

	timeNow := time.Now()
	// update
	userUpdate := domain.User{
		Uuid:      req.Uuid,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Email:     user.Email,
		Username:  req.Username,
		Password:  req.Password,
		UpdatedAt: &timeNow,
	}

	err = u.userRepository.Update(ctx, &userUpdate)
	if err != nil {
		return userDto.UserUpdateResponse{}, err
	}

	return userDto.UserUpdateResponse{
		Uuid:          user.Uuid,
		FullName:      userUpdate.FullName,
		Phone:         userUpdate.Phone,
		Email:         user.Email,
		Username:      userUpdate.Username,
		EmailVerifyAt: user.EmailVerifyAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     userUpdate.UpdatedAt,
	}, err
}
