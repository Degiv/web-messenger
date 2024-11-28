package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"messenger/internals/domain"
)

type Users struct {
	DB *sqlx.DB
}

func NewUsers(DB *sqlx.DB) *Users {
	return &Users{
		DB: DB,
	}
}

func (u *Users) GetUsersByIDs(userIDs []int64) ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	err := u.DB.Select(&users, "SELECT * FROM users WHERE id = ANY ($1)", pq.Array(userIDs))
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *Users) GetUserByID(userID int64) (*domain.User, error) {
	const query = `SELECT * FROM users WHERE id = $1`
	var user *domain.User
	err := u.DB.Get(user, query, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *Users) GetUserByUsername(username string) (*domain.User, error) {
	const query = `SELECT * FROM users WHERE username = $1`
	var user *domain.User
	err := u.DB.Get(user, query, username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
