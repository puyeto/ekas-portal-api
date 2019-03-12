package models

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
