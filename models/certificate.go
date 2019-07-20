package models

import "github.com/go-ozzo/ozzo-validation"

// Artist represents an artist record.
type Certificates struct {
	ID  int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Validate validates the Artist fields.
func (m Certificates) ValidateCertificates() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 120)),
	)
}
