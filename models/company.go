package models

import validation "github.com/go-ozzo/ozzo-validation"

// Companies ...
type Companies struct {
	CompanyID          int    `json:"company_id" db:"pk,company_id"`
	CompanyName        string `json:"company_name" db:"company_name"`
	CompanyContacts    string `json:"company_contacts"`
	CompanyContactName string `json:"company_contact_name"`
	CompanyEmail       string `json:"company_email"`
	CompanyLocation    string `json:"company_location"`
	UserID             int32  `json:"user_id,omitempty"`
}

// ValidateCompanies validates user data fields.
func (m Companies) ValidateCompanies() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CompanyName, validation.Required),
		validation.Field(&m.CompanyContactName, validation.Required),
		validation.Field(&m.CompanyContacts, validation.Required),
		validation.Field(&m.CompanyEmail, validation.Required),
		validation.Field(&m.CompanyLocation, validation.Required),
	)
}
