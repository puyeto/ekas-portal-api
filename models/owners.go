package models

import validation "github.com/go-ozzo/ozzo-validation"

// VehicleOwner ...
type VehicleOwner struct {
	OwnerID    uint32 `json:"owner_id" db:"pk,owner_id"`
	UserID     uint32 `json:"user_id" db:"user_id"`
	OwnerIDNo  string `json:"owner_id_no" db:"owner_id_no"`
	OwnerName  string `json:"owner_name" db:"owner_name"`
	OwnerEmail string `json:"owner_email" db:"owner_email"`
	OwnerPhone string `json:"owner_phone" db:"owner_phone"`
}

// ValidateVehicleOwner validates fields.
func (v VehicleOwner) ValidateVehicleOwner() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.UserID, validation.Required),
		validation.Field(&v.OwnerIDNo, validation.Required),
		validation.Field(&v.OwnerName, validation.Required),
	)
}
