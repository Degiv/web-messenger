package domain

import "time"

type User struct {
	ID           int64      `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"password_hash"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at"`
	LastLogin    *time.Time `db:"last_login" json:"last_login"`
}
