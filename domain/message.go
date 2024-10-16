package domain

import "time"

type Message struct {
	ID           int64     `db:"id" json:"id"`
	ConferenceID int64     `db:"conference_id" json:"conference_id"`
	UserID       int64     `db:"user_id" json:"user_id"`
	Text         string    `db:"text" json:"text"`
	SentAt       time.Time `db:"sent_at" json:"sent_at"`
}
