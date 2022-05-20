package models

import validation "github.com/go-ozzo/ozzo-validation"

// Saccos ...
type Saccos struct {
	ID      int32  `json:"id" db:"pk,id"`
	Name    string `json:"name" db:"name"`
	Alias   string `json:"alias" db:"short_name"`
	Address string `json:"address" db:"address"`
}

// Validate validates user data fields.
func (s Saccos) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Alias, validation.Required),
	)
}
