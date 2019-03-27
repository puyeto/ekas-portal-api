package models

import validation "github.com/go-ozzo/ozzo-validation"

// TrackingServerAuth represents an trackingServer record.
type TrackingServerAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ValidateTrackingServerLogin validates the TrackingServerAuth fields.
func (m TrackingServerAuth) ValidateTrackingServerLogin() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, validation.Length(0, 120)),
		validation.Field(&m.Password, validation.Required, validation.Length(0, 120)),
	)
}

// AddDevice ...
type AddDevice struct {
	UserData   UserData         `json:"user_data"`
	DeviceData AddDeviceDetails `json:"device_data"`
}

// AddDeviceDetails ...
type AddDeviceDetails struct {
	Name               string `json:"name"`
	Imei               string `json:"imei"`
	IconID             string `json:"icon_id"`
	FuelMeasurementID  string `json:"fuel_measurement_id"`
	TailLength         string `json:"tail_length"`
	MinMovingSpeed     string `json:"min_moving_speed"`
	MinFuelFillings    string `json:"min_fuel_fillings"`
	MinFuelThefts      string `json:"min_fuel_thefts"`
	PlateNumber        string `json:"plate_number"`
	Vin                string `json:"vin"`
	DeviceModel        string `json:"device_model"`
	RegistrationNumber string `json:"registration_number"`
	ObjectOwner        string `json:"object_owner"`
}

// ValidateAddDevices validates addition of devices fields.
func (m AddDeviceDetails) ValidateAddDevices() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 120)),
		validation.Field(&m.Imei, validation.Required, validation.Length(0, 120)),
		validation.Field(&m.IconID, validation.Required),
		validation.Field(&m.FuelMeasurementID, validation.Required),
		validation.Field(&m.TailLength, validation.Required),
		validation.Field(&m.MinMovingSpeed, validation.Required),
		validation.Field(&m.MinFuelFillings, validation.Required),
		validation.Field(&m.MinFuelThefts, validation.Required),
	)
}

type AddServices struct {
	UserData                  UserData                  `json:"device_user_data"`
	AddTrackingServiceDetails AddTrackingServiceDetails `json:"add_service"`
}

// AddTrackingServiceDetails ...
type AddTrackingServiceDetails struct {
	Name                 string `json:"name"`
	DeviceID             int32  `json:"device_id"`
	ExpirationBy         string `json:"expiration_by"`
	Interval             string `json:"interval"`
	LastService          string `json:"last_service"`
	TriggerEventLeft     string `json:"trigger_event_left"`
	RenewAfterExpiration int8   `json:"renew_after_expiration"`
	Email                string `json:"email"`
	MobilePhone          string `json:"mobile_phone"`
}

// ValidateAddTrackingServiceDetails validates addition of service fields.
func (m AddTrackingServiceDetails) ValidateAddTrackingServiceDetails() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 120)),
		validation.Field(&m.DeviceID, validation.Required),
		validation.Field(&m.ExpirationBy, validation.Required),
		validation.Field(&m.LastService, validation.Required),
		validation.Field(&m.TriggerEventLeft, validation.Required),
		validation.Field(&m.RenewAfterExpiration, validation.NotNil),
		validation.Field(&m.Email, validation.Required),
		validation.Field(&m.MobilePhone, validation.Required),
	)
}

// TrackingServiceTypes ...
type TrackingServiceTypes struct {
	Odometer    string `json:"odometer"`     // Odometer
	EngineHours string `json:"engine_hours"` // Engine hours
	Days        string `json:"days"`         // Days
}

// TrackingServerService ...
type TrackingServerService struct {
	ID                   string      `json:"id"`
	DeviceID             string      `json:"device_id"`
	Name                 string      `json:"name"`
	ExpirationBy         string      `json:"expiration_by"`
	Interval             string      `json:"interval"`
	LastService          string      `json:"last_service"`
	TriggerEventLeft     string      `json:"trigger_event_left"`
	RenewAfterExpiration string      `json:"renew_after_expiration"`
	Expires              string      `json:"expires"`
	ExpiresDate          string      `json:"expires_date,omitempty"`
	Remind               string      `json:"remind"`
	RemindDate           string      `json:"remind_date,omitempty"`
	EventSent            string      `json:"event_sent"`
	Expired              string      `json:"expired"`
	Email                string      `json:"email,omitempty"`
	MobilePhone          string      `json:"mobile_phone,omitempty"`
	Device               interface{} `json:"devices"`
}
