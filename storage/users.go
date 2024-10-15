package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"messenger/domain"
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

func (u *Users) CreateRelationToConference(userID int64, conferenceID int64) error {
	const query = `INSERT INTO usersConferencesRelation (user_id, conference_id) VALUES ($1, $2)`
	_, err := u.DB.Exec(query, userID, conferenceID)
	return err
}
