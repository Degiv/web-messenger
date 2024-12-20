package auth

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Degiv/web-messenger/internals/domain"
	"github.com/Degiv/web-messenger/pkg/passwordhashing"
)

var (
	ErrWrongPassword = errors.New("wrong password")
	ErrNoSuchUser    = errors.New("no such user")
)

type Users interface {
	GetUserByUsername(username string) (*domain.User, error)
}

type Service struct {
	Users Users
}

func NewAuthService(users Users) *Service {
	return &Service{
		Users: users,
	}
}

func (s *Service) AuthorizeUser(username string, password string) (int64, error) {
	user, err := s.Users.GetUserByUsername(username)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNoSuchUser
	}

	if err != nil {
		return 0, fmt.Errorf("get user by username: %w", err)
	}

	if !passwordhashing.VerifyPassword(password, user.PasswordHash) {
		return 0, ErrWrongPassword
	}

	return user.ID, nil
}
