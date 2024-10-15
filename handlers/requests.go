package handlers

import "time"

type createConferenceRequest struct {
	usersIDs  []int64
	name      string
	createdBy int64
	createdAt time.Time
}
