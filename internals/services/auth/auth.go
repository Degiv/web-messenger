package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"messenger/internals/domain"
	"messenger/pkg/passwordhashing"
)

var (
	ErrWrongPassword = errors.New("Wrong password")
	ErrNoSuchUser    = errors.New("No such user")
)

type Users interface {
	GetUserByUsername(username string) (*domain.User, error)
}

type AuthService struct {
	Users Users
}

func NewAuthService(users Users) *AuthService {
	return &AuthService{
		Users: users,
	}
}

func (service *AuthService) AuthorizeUser(username string, password string) (int64, error) {
	user, err := service.Users.GetUserByUsername(username)
	if err == sql.ErrNoRows {
		return 0, ErrNoSuchUser
	}

	if err != nil {
		return 0, fmt.Errorf("Get user by username: %w", err)
	}

	if !passwordhashing.VerifyPassword(password, user.PasswordHash) {
		return 0, ErrWrongPassword
	}

	return user.ID, nil
}
