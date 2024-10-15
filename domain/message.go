package domain

import "time"

type Message struct {
	ID           int64     `db:"id"`
	ConferenceID int64     `db:"conference_id"`
	UserID       int64     `db:"user_id"`
	Text         string    `db:"text"`
	SentAt       time.Time `db:"sent_at"`
}
