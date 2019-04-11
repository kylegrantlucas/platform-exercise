package models

import "time"

type Session struct {
	UUID      string    `json:"uuid,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	UserUUID  string    `json:"user_uuid,omitempty"`
}
