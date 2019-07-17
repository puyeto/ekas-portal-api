package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Devices ..,
type Devices struct {
	ID                 int32     `json:"id" db:"id"`
	DeviceID           int32     `json:"device_id" db:"device_id"`
	DeviceName         string    `json:"device_name,omitempty" db:"device_name"`
	DeviceSerialNo     string    `json:"device_serial_no" db:"device_serial_no"`
	DeviceModelNo      string    `json:"device_model_no" db:"device_model_no"`
	DeviceManufacturer string    `json:"device_manufacturer" db:"device_manufacturer"`
	Configured         int8      `json:"configured" db:"configured"`
	Status             int8      `json:"status" db:"status"`
	Note               string    `json:"note" db:"note"`
	CreatedOn          time.Time `json:"created_on" db:"created_on"`
}

// ValidateDevices validates user data fields.
func (m Devices) ValidateDevices() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeviceName, validation.Required),
	)
}

// DeviceConfiguration ...
type DeviceConfiguration struct {
	DeviceID   int32  `json:"device_id"`
	DeviceName string `json:"device_name,omitempty" db:"device_name"`
	SIMImei    string `json:"sim_imei" db:"sim_imei"`
}
