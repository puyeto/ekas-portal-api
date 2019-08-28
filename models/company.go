package models

import validation "github.com/go-ozzo/ozzo-validation"

// Companies ...
type Companies struct {
	CompanyID          int32  `json:"company_id" db:"pk,company_id"`
	CompanyName        string `json:"company_name" db:"company_name"`
	CompanyContacts    string `json:"company_contacts"`
	CompanyContactName string `json:"company_contact_name"`
	CompanyEmail       string `json:"company_email"`
	CompanyLocation    string `json:"company_location"`
	UserID             int32  `json:"user_id" db:"updated_by"`
	ContactID          int32  `json:"contact_id,omitempty" db:"contact_id"`
	CompanyPhone       string `json:"company_phone"`
	BusinessRegNo      string `json:"business_reg_no"`
	User               int32  `json:"user,omitempty" db:"user"`
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
