package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Devices ..,
type Devices struct {
	ID                 int32        `json:"id" db:"id"`
	DeviceID           int64        `json:"device_id" db:"device_id"`
	VehicleID          uint32       `json:"vehicle_id" db:"vehicle_id"`
	DeviceName         string       `json:"device_name" db:"device_name"`
	CompanyID          int          `json:"company_id" db:"company_id"`
	CompanyName        string       `json:"company_name,omitempty" db:"company_name"`
	DeviceSerialNo     string       `json:"device_serial_no,omitempty" db:"device_serial_no"`
	DeviceModelNo      string       `json:"device_model_no,omitempty" db:"device_model"`
	MotherboardNO      string       `json:"motherboard_no,omitempty" db:"motherboard_no"`
	Technician         string       `json:"technician,omitempty" db:"technician"`
	DeviceManufacturer string       `json:"device_manufacturer,omitempty" db:"device_manufacturer"`
	Configured         int8         `json:"configured,omitempty" db:"configured"`
	Status             int8         `json:"status,omitempty" db:"status"`
	SimSerialNo        string       `json:"sim_serial_no,omitempty" db:"sim_serial_no"`
	SimNumber          string       `json:"sim_number,omitempty" db:"sim_number"`
	Note               string       `json:"note,omitempty" db:"note"`
	CreatedOn          time.Time    `json:"created_on,omitempty" db:"created_on"`
	VehicleRegNno      string       `json:"vehicle_reg_no,omitempty" db:"vehicle_reg_no"`
	ChassisNo          string       `json:"chassis_no,omitempty" db:"chassis_no"`
	MakeType           string       `json:"make_type,omitempty" db:"make_type"`
	Model              string       `json:"model,omitempty" db:"model"`
	ModelYear          int          `json:"model_year,omitempty" db:"model_year"`
	Positions          []DeviceData `json:"positions,omitempty"`
}

// NewDevice ...
func NewDevice(did int64, dname, dsno, sno, mno, tech string) *Devices {
	nd := &Devices{
		DeviceID:           did,
		DeviceName:         dname,
		DeviceSerialNo:     dsno,
		DeviceModelNo:      "Ekas FK 2017",
		DeviceManufacturer: "Ekkas Technologies",
		SimNumber:          sno,
		MotherboardNO:      mno,
		Technician:         tech,
	}

	return nd
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
	ConfigID            int32     `json:"conf_id" db:"pk,conf_id"`
	DeviceID            int32     `json:"device_id" db:"device_id"`
	DeviceName          string    `json:"device_name,omitempty" db:"dveame"`
	ChassisNo           string    `json:"chassis_no,omitempty" db:"chis_no"`
	MakeType            string    `json:"make_type" db:"make_type"`
	DeviceType          string    `json:"device_type" db:"device_type"`
	SerialNo            string    `json:"serial_no" db:"serial_no"`
	SIMImei             string    `json:"sim_imei" db:"sim_imei"`
	VehicleID           string    `json:"vehicle_id" db:"vehicle_id"`
	CreatedOn           time.Time `json:"created_on" db:"created_on"`
	ConfigurationStatus int8      `json:"status" db:"status"`
	DeviceStatus        int8      `json:"device_status" db:"device_status"`
	DeviceReason        string    `json:"device_reason" db:"reason"`
	PresetSpeed         string    `json:"preset_speed" db:"preset_speed"`
	SetFrequency        string    `json:"set_frequency" db:"set_frequency"`
	FittingDate         string    `json:"fitting_date" db:"fitting_date"`
	ExpiryDate          string    `json:"expiry_date" db:"expiry_date"`
	FittingCenter       string    `json:"fitting_center" db:"fitting_center"`
	Certificate         string    `json:"certificate" db:"certificate"`
	EmailAddress        string    `json:"agent_email_address" db:"email_address"`
	AgentPhone          string    `json:"agent_phone" db:"agent_phone"`
	AgentLocation       string    `json:"agent_location" db:"agent_location"`
	OwnerName           string    `json:"owner_name" db:"owner_name"`
	OwnerPhoneNumber    string    `json:"owner_phone_number" db:"owner_phone_number"`
}
