package models

import (
	"encoding/json"
	"time"
)

type User struct {
	UUID      string    `json:"uuid,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// MarshalJSON is a custom marshaler for the User struct that strips the password before marshaling
func (u *User) MarshalJSON() ([]byte, error) {
	userWithoutPassword := *u
	userWithoutPassword.Password = ""
	return json.Marshal(userWithoutPassword)
}
