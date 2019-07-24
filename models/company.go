package models

import validation "github.com/go-ozzo/ozzo-validation"

// Companies ...
type Companies struct {
	CompanyID       int    `json:"company_id" db:"pk,company_id"`
	CompanyName     string `json:"company_name" db:"company_name"`
	CompanyContacts string `json:"company_contacts,omitempty"`
	UserID          int32  `json:"user_id,omitempty"`
}

// ValidateDevices validates user data fields.
func (m Companies) ValidateCompanies() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CompanyName, validation.Required),
	)
}
