package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Simcards ...
type Simcards struct {
	ID             uint32    `json:"id" db:"pk,id"`
	Identifier     string    `json:"identifier" db:"identifier"`
	ICCID          string    `json:"iccid" db:"iccid"`
	Network        string    `json:"network" db:"network"`
	SIMType        string    `json:"sim_type" db:"sim_type"`
	Managed        int8      `json:"managed" db:"managed"`
	Description    string    `json:"description" db:"description"`
	AirtimeBalance float32   `json:"airtime_balance" db:"airtime_balance"`
	DataBalance    float32   `json:"data_balance" db:"data_balance"`
	SMSBalance     float32   `json:"sms_balance" db:"sms_balance"`
	AddedAt        time.Time `json:"added_at" db:"added_at"`
	RuleCount      int8      `json:"rule_count" db:"rule_count"`
}

// Validate fields.
func (s Simcards) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Identifier, validation.Required),
	)
}

// SimcardStats ...
type SimcardStats struct {
	TotalCount     int `json:"total_count"`
	ManagedCount   int `json:"managed_count"`
	UnManagedCount int `json:"un_managed_count"`
}
