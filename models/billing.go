package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Discounts ...
type Discounts struct {
	DiscountID      int32   `json:"discount_id,omitempty" db:"pk,discount_id"`
	DiscountName    string  `json:"discount_name" db:"discount_name"`
	DiscountValue   float32 `json:"discount_value" db:"discount_value"`
	DiscountValueAs string  `json:"discount_value_as" db:"discount_value_as"`
}

// ValidateDiscount ...
func (d Discounts) ValidateDiscount() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.DiscountName, validation.Required),
		validation.Field(&d.DiscountValue, validation.Required),
	)
}

// Pricings ...
type Pricings struct {
	PricingID       int32   `json:"pricing_id,omitempty" db:"pk,pricing_id"`
	PricingName     string  `json:"pricing_name" db:"pricing_name"`
	PricingAmount   float32 `json:"pricing_amount" db:"pricing_amount"`
	PricingDuration string  `json:"pricing_duration" db:"pricing_duration"`
}

// ValidateDiscount ...
func (p Pricings) ValidatePricing() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.PricingName, validation.Required),
		validation.Field(&p.PricingAmount, validation.Required),
		validation.Field(&p.PricingDuration, validation.Required),
	)
}

// Invoices ...
type Invoices struct {
	InvoiceID        int32     `json:"invoice_id,omitempty" db:"pk,invoice_id"`
	VehicleID        int32     `json:"vehicle_id" db:"vehicle_id"`
	SettingID        int32     `json:"vehicle_id" db:"setting_id"`
	PricingID        int32     `json:"pricing_id" db:"pricing_id"`
	DiscountID       int32     `json:"discount_id" db:"discount_id"`
	InvoiceName      string    `json:"invoice_name,omitempty" db:"invoice_name"`
	InvoiceAmount    float32   `json:"invoice_amount,omitempty" db:"invoice_amount"`
	InvoiceTotal     float32   `json:"invoice_total,omitempty" db:"invoice_total"`
	InvoicingDate    time.Time `json:"invoicing_date,omitempty" db:"invoicing_date"`
	InvoiceStartDate time.Time `json:"invoice_start_date,omitempty" db:"invoice_start_date"`
	InvoiceEndDate   time.Time `json:"invoice_end_date,omitempty" db:"invoice_end_date"`
	InvoiceStatus   time.Time `json:"invoice_status,omitempty" db:"invoice_status"`
	InvoicePaymentStatus   time.Time `json:"invoice_payment_status,omitempty" db:"invoice_payment_status"`
}

// ValidateInvoice ...
func (in Invoices) ValidateInvoice() error {
	return validation.ValidateStruct(&in,
		validation.Field(&in.VehicleID, validation.Required),
		validation.Field(&in.PricingID, validation.Required),
		validation.Field(&in.DiscountID, validation.Required),
		validation.Field(&in.InvoiceAmount, validation.Required),
		validation.Field(&in.InvoiceTotal, validation.Required),
		validation.Field(&in.InvoicingDate, validation.Required),
		validation.Field(&in.InvoiceStartDate, validation.Required),
		validation.Field(&in.InvoiceEndDate, validation.Required),
	)
}
