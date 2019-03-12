package models

import validation "github.com/go-ozzo/ozzo-validation"

// TrackingServerAuth represents an trackingServer record.
type TrackingServerAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate validates the TrackingServerAuth fields.
func (m TrackingServerAuth) ValidateTrackingServerLogin() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, validation.Length(0, 120)),
		validation.Field(&m.Password, validation.Required, validation.Length(0, 120)),
	)
}
