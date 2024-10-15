package services

import (
	"fmt"
	"messenger/domain"
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
	GetUsersByIDs(usersIDs []int64) ([]*domain.User, error)
}

type MessengerService struct {
	Users       Users
	Messages    Messages
	Conferences Conferences
}

func NewMessenger(users Users, messages Messages, conferences Conferences) *MessengerService {
	return &MessengerService{
		Users:       users,
		Messages:    messages,
		Conferences: conferences,
	}
}

func (service *MessengerService) GetConferenceMessages(userID int64) ([]*domain.Message, error) {
	return service.Messages.GetMessagesByConference(userID)
}

func (service *MessengerService) PostToConference(message *domain.Message) error {
	return service.Messages.InsertMessage(message)
}

func (service *MessengerService) GetConferencesByUser(userID int64) ([]*domain.Conference, error) {
	return service.Conferences.GetConferencesByUser(userID)
}

func (service *MessengerService) CreateConference(usersIDs []int64, name string, createdBy int64, createdAt time.Time) error {
	users, err := service.Users.GetUsersByIDs(usersIDs)
	if err != nil {
		return fmt.Errorf("validating users: %w", err)
	}

	if len(users) != len(usersIDs) {
		return fmt.Errorf("validating users: some users doesn't exist")
	}

	conferenceID, err := service.Conferences.CreateConference(name, createdBy, createdAt)
	if err != nil {
		return fmt.Errorf("creating conference: %w", err)
	}

	for _, userID := range usersIDs {
		err = service.Conferences.CreateConferenceMember(userID, conferenceID, createdAt)
		if err != nil {
			return fmt.Errorf("creating usersConferencesRelation: %w", err)
		}
	}

	return nil
}
