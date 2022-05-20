package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Certificates represents an certificate record.
type Certificates struct {
	ID         int       `json:"id" db:"id"`
	CertSerial string    `json:"cert_serial"`
	CertNo     string    `json:"cert_no"`
	CompanyID  int       `json:"company_id"`
	FitterID   int       `json:"fitter_id"`
	Company    string    `json:"company"` // agent / fitting center
	Fitter     string    `json:"fitter"`  // agent technician name
	IssuedOn   time.Time `json:"issued_on"`
}

// ValidateCertificates validates the certificate fields.
func (m Certificates) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CertNo, validation.Required, validation.Length(0, 120)),
		validation.Field(&m.CompanyID, validation.Required),
		validation.Field(&m.FitterID, validation.Required),
	)
}
