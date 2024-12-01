package domain

import "time"

type Conference struct {
	ID          int64      `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	CreatedBy   int64      `db:"created_by" json:"created_by"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	LastMessage *int64     `db:"last_message" json:"last_message"`
}
