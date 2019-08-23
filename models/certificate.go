package models

import validation "github.com/go-ozzo/ozzo-validation"

// Certificates represents an certificate record.
type Certificates struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// ValidateCertificates validates the certificate fields.
func (m Certificates) ValidateCertificates() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 120)),
	)
}
