package storage

import (
	"github.com/jmoiron/sqlx"
	"messenger/internals/domain"
	"time"
)

type Conferences struct {
	DB *sqlx.DB
}

func NewConferences(DB *sqlx.DB) *Conferences {
	return &Conferences{
		DB: DB,
	}
}

func (c *Conferences) GetConferencesByUser(userID int64) ([]*domain.Conference, error) {
	conferences := make([]*domain.Conference, 0)
	const query = `SELECT * FROM conferences AS c 
    JOIN conference_members AS cm ON c.id = cm.conference_id
    WHERE cm.user_id = $1`
	err := c.DB.Select(&conferences, query, userID)
	if err != nil {
		return nil, err
	}

	return conferences, nil
}

func (c *Conferences) CreateConference(name string, createdBy int64, createdAt time.Time) (int64, error) {
	const query = `INSERT INTO conferences (name, created_at, created_by) VALUES ($1, $2, $3) RETURNING id;`
	row := c.DB.QueryRowx(query, &name, &createdAt, &createdBy)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (c *Conferences) CreateConferenceMember(userID int64, conferenceID int64, joinedAt time.Time) error {
	const query = `INSERT INTO conference_members (user_	id, conference_id, joined_at) VALUES ($1, $2, $3)`
	_, err := c.DB.Exec(query, &userID, &conferenceID, &joinedAt)
	return err
}
