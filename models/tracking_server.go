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
	Name              string `json:"name"`
	Imei              string `json:"imei"`
	IconID            string `json:"icon_id"`
	FuelMeasurementID string `json:"fuel_measurement_id"`
	TailLength        string `json:"tail_length"`
	MinMovingSpeed    string `json:"min_moving_speed"`
	MinFuelFillings   string `json:"min_fuel_fillings"`
	MinFuelThefts     string `json:"min_fuel_thefts"`
	PlateNumber       string `json:"plate_number"`
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
