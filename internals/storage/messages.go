package storage

import (
	"github.com/jmoiron/sqlx"

	"github.com/Degiv/web-messenger/internals/domain"
)

type Messages struct {
	DB *sqlx.DB
}

func NewMessages(db *sqlx.DB) *Messages {
	return &Messages{
		DB: db,
	}
}

func (m *Messages) GetMessagesByConference(id int64) ([]*domain.Message, error) {
	messages := make([]*domain.Message, 0)
	err := m.DB.Select(&messages, "SELECT * FROM messages WHERE conference_id = $1", id)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *Messages) InsertMessage(message *domain.Message) error {
	_, err := m.DB.Exec("INSERT INTO messages (conference_id, user_id, text, sent_at) VALUES ($1, $2, $3, $4)",
		message.ConferenceID,
		message.UserID,
		message.Text,
		message.SentAt)

	return err
}
