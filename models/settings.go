package models

import validation "github.com/go-ozzo/ozzo-validation"

// Settings ..
type Settings struct {
	SettingID         int     `json:"setting_id" db:"pk,setting_id"`
	CertificatePrefix string  `json:"certificate_prefix" db:"certificate_prefix"`
	DevicePrefix      int     `json:"device_prefix" db:"device_prefix"`
	BillingPeriod     int8    `json:"billing_period" db:"billing_period"`
	MonthlyAmount     float32 `json:"monthly_amount" db:"monthly_amount"`
	InfoMessage string `json:"info_message" db:"info_message"`
}

// ValidateSettings ...
func (s Settings) ValidateSettings() error {
	return validation.ValidateStruct(&s) // validation.Field(&s.CompanyName, validation.Required, validation.Length(3, 120)),

}

// LicenseKeys List keys
type LicenseKeys struct {
	KeyString       string `json:"key_string" db:"pk,key_string"`
	AssignTo        int    `json:"assign_to"`
	Status          int8   `json:"status"`
	CompanyID       int    `json:"company_id,omitempty"`
	CompanyName     string `json:"company_name,omitempty"`
	CompanyContacts string `json:"company_contacts,omitempty"`
}

// ValidateLicenseKeys ...
func (s LicenseKeys) ValidateLicenseKeys() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.KeyString, validation.Required),
		validation.Field(&s.AssignTo, validation.Required),
	)
}

// GenKeys generate keys
type GenKeys struct {
	Number   int    `json:"number"`
	Type     string `json:"type"`
	AssignTo int    `json:"assign_to"`
}

// ValidateGenKeys ...
func (s GenKeys) ValidateGenKeys() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Number, validation.Required),
		validation.Field(&s.Type, validation.Required),
		validation.Field(&s.AssignTo, validation.Required),
	)
}
