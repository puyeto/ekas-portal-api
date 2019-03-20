package models

import validation "github.com/go-ozzo/ozzo-validation"

// Identity ..
type Identity interface {
	GetID() string
	GetName() string
}

// UserData ...
type UserData struct {
	Lang     string `json:"lang"`
	UserHash string `json:"user_api_hash"`
}

// ValidateUserData validates user data fields.
func (m UserData) ValidateUserData() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Lang, validation.Required),
		validation.Field(&m.UserHash, validation.Required),
	)
}

// User ...
type User struct {
	ID   string
	Name string
}

// GetID ...
func (u User) GetID() string {
	return u.ID
}

// GetName ...
func (u User) GetName() string {
	return u.Name
}
