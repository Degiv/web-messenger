package handlers

import "time"

type createConferenceRequest struct {
	UsersIDs  []int64   `json:"users_ids"`
	Name      string    `json:"name"`
	CreatedBy int64     `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
