package models

import validation "github.com/go-ozzo/ozzo-validation"

// Identity ..
type Settings struct {
	SettingID   int    `json:"setting_id" db:"pk,setting_id"`
	CompanyName string `json:"company_name" db:"company_name"`
}

// ValidateSettings ...
func (s Settings) ValidateSettings() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.CompanyName, validation.Required, validation.Length(6, 120)),
	)
}

// LicenseKeys List keys
type LicenseKeys struct {
	KeyString string `json:"key_string" db:"pk,key_string"`
	AssignTo  int    `json:"assign_to"`
	Status    int    `json:"status,omitempty"`
}

// ValidateGenKeys ...
func (s LicenseKeys) ValidateLicenseKeys() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.KeyString, validation.Required),
		validation.Field(&s.AssignTo, validation.Required),
	)
}

// GenKeys generate keys
type GenKeys struct {
	Number int    `json:"number"`
	Type   string `json:"type"`
}

// ValidateGenKeys ...
func (s GenKeys) ValidateGenKeys() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Number, validation.Required),
		validation.Field(&s.Type, validation.Required),
	)
}
