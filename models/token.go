package models

import "time"

type Token struct {
	Value     string    `json:"token,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	UserUUID  string    `json:"user_uuid,omitempty"`
}
