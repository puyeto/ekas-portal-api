package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Devices ..,
type Devices struct {
	ID                 int32     `json:"id" db:"id"`
	DeviceID           int32     `json:"device_id" db:"device_id"`
	DeviceName         string    `json:"device_name" db:"device_name"`
	DeviceSerialNo     string    `json:"device_serial_no,omitempty" db:"device_serial_no"`
	DeviceModelNo      string    `json:"device_model_no,omitempty" db:"device_model"`
	DeviceManufacturer string    `json:"device_manufacturer,omitempty" db:"device_manufacturer"`
	Configured         int8      `json:"configured,omitempty" db:"configured"`
	Status             int8      `json:"status,omitempty" db:"status"`
	Note               string    `json:"note,omitempty" db:"note"`
	CreatedOn          time.Time `json:"created_on,omitempty" db:"created_on"`
}

// ValidateDevices validates user data fields.
func (m Devices) ValidateDevices() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeviceID, validation.Required),
		validation.Field(&m.DeviceSerialNo, validation.Required),
		validation.Field(&m.DeviceModelNo, validation.Required),
	)
}

// DeviceConfiguration ...
type DeviceConfiguration struct {
	DeviceID   int32  `json:"device_id"`
	DeviceName string `json:"device_name,omitempty" db:"device_name"`
	SIMImei    string `json:"sim_imei" db:"sim_imei"`
}
