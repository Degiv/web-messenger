package messenger

import (
	"fmt"
	"messenger/internals/domain"
	"messenger/pkg/passwordhashing"
	"time"
)

type Conferences interface {
	GetConferencesByUser(userID int64) ([]*domain.Conference, error)
	CreateConference(name string, createdBy int64, createdAt time.Time) (int64, error)
	CreateConferenceMember(userID int64, conferenceID int64, joinedAt time.Time) error
}

type Messages interface {
	GetMessagesByConference(id int64) ([]*domain.Message, error)
	InsertMessage(message *domain.Message) error
}

type Users interface {
	CreateUser(username string, email string, passwordHash string) (int64, error)
	GetUsersByIDs(usersIDs []int64) ([]*domain.User, error)
	GetUserByID(userID int64) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
}

type Service struct {
	Users       Users
	Messages    Messages
	Conferences Conferences
}

func NewMessenger(users Users, messages Messages, conferences Conferences) *Service {
	return &Service{
		Users:       users,
		Messages:    messages,
		Conferences: conferences,
	}
}

func (s *Service) GetUserByID(userID int64) (*domain.User, error) {
	user, err := s.Users.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetConferenceMessages(userID int64) ([]*domain.Message, error) {
	return s.Messages.GetMessagesByConference(userID)
}

func (s *Service) PostToConference(message *domain.Message) error {
	return s.Messages.InsertMessage(message)
}

func (s *Service) GetConferencesByUser(userID int64) ([]*domain.Conference, error) {
	return s.Conferences.GetConferencesByUser(userID)
}

func (s *Service) CreateConference(usersIDs []int64, name string, createdBy int64, createdAt time.Time) error {
	users, err := s.Users.GetUsersByIDs(usersIDs)
	if err != nil {
		return fmt.Errorf("validating users: %w", err)
	}

	if len(users) != len(usersIDs) {
		return fmt.Errorf("validating users: some users doesn't exist")
	}

	conferenceID, err := s.Conferences.CreateConference(name, createdBy, createdAt)
	if err != nil {
		return fmt.Errorf("creating conference: %w", err)
	}

	for _, userID := range usersIDs {
		err = s.Conferences.CreateConferenceMember(userID, conferenceID, createdAt)
		if err != nil {
			return fmt.Errorf("creating usersConferencesRelation: %w", err)
		}
	}

	return nil
}

func (s *Service) RegisterUser(username string, email string, password string) (int64, error) {
	passwordHash, err := passwordhashing.HashPassword(password)
	if err != nil {
		return 0, err
	}

	id, err := s.Users.CreateUser(username, email, passwordHash)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) VerifyConferenceMember(userID int64, conferenceID int64) (bool, error) {
	conferences, err := s.Conferences.GetConferencesByUser(userID)
	if err != nil {
		return false, err
	}

	ok := false
	for _, c := range conferences {
		if c.ID == conferenceID {
			ok = true
			break
		}
	}

	return ok, nil
}
