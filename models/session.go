package models

import "time"

type Session struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	UserUUID  string    `json:"user_uuid,omitempty"`
}
