package storage

import (
	"github.com/jmoiron/sqlx"
	"messenger/internals/domain"
)

type MessagesStorage struct {
	DB *sqlx.DB
}

func NewMessages(DB *sqlx.DB) *MessagesStorage {
	return &MessagesStorage{
		DB: DB,
	}
}

func (storage *MessagesStorage) GetMessagesByConference(id int64) ([]*domain.Message, error) {
	messages := make([]*domain.Message, 0)
	err := storage.DB.Select(&messages, "SELECT * FROM messages WHERE conference_id = $1", id)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (storage *MessagesStorage) InsertMessage(message *domain.Message) error {
	_, err := storage.DB.Exec("INSERT INTO messages (id, conference_id, user_id, text, sent_at) VALUES ($1, $2, $3, $4, $5)",
		message.ID,
		message.ConferenceID,
		message.UserID,
		message.Text,
		message.SentAt)

	return err
}
