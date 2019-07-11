package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Devices ..,
type Devices struct {
	DeviceID   int32  `json:"device_id"`
	DeviceName string `json:"device_name,omitempty" db:"device_name"`
	SIMImei    string `json:"sim_imei" db:"sim_imei"`
}

// ValidateDevices validates user data fields.
func (m Devices) ValidateDevices() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeviceName, validation.Required),
	)
}
