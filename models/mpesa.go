package models

import (
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation"
)

// AuthResponse ...
type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

// C2B is a model
type C2B struct {
	ShortCode     string
	CommandID     string
	Amount        string
	Msisdn        string
	BillRefNumber string
}

// C2BSimulateResponse ...
type C2BSimulateResponse struct {
	ConversationID          string
	OriginatorCoversationID string
	ResponseDescription     string
	Status                  int `json:"status"`
}

//ValidateC2B validate C2B Transaction
func (u C2B) ValidateC2B() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ShortCode, validation.Required),
		validation.Field(&u.CommandID, validation.Required),
		validation.Field(&u.Amount, validation.Required),
		validation.Field(&u.Msisdn, validation.Required),
		validation.Field(&u.BillRefNumber, validation.Required),
	)
}

// B2C is a model
type B2C struct {
	InitiatorName      string
	SecurityCredential string
	CommandID          string
	Amount             string
	PartyA             string
	PartyB             string
	Remarks            string
	QueueTimeOutURL    string
	ResultURL          string
	Occassion          string
}

// B2B is a model
type B2B struct {
	Initiator              string
	SecurityCredential     string
	CommandID              string
	SenderIdentifierType   string
	RecieverIdentifierType string
	Amount                 string
	PartyA                 string
	PartyB                 string
	Remarks                string
	AccountReference       string
	QueueTimeOutURL        string
	ResultURL              string
}

// Express is a model
type Express struct {
	BusinessShortCode string
	Password          string
	Timestamp         string
	TransactionType   string
	Amount            string
	PartyA            string
	PartyB            string
	PhoneNumber       string
	CallBackURL       string
	AccountReference  string
	TransactionDesc   string
}

// Status check transaction status
type Status struct {
	BusinessShortCode string
	Password          string
	Timestamp         string
	CheckoutRequestID string
}

// Reversal is a model
type Reversal struct {
	Initiator              string
	SecurityCredential     string
	CommandID              string
	TransactionID          string
	Amount                 string
	ReceiverParty          string
	ReceiverIdentifierType string
	QueueTimeOutURL        string
	ResultURL              string
	Remarks                string
	Occassion              string
}

// BalanceInquiry is a model
type BalanceInquiry struct {
	Initiator          string
	SecurityCredential string
	CommandID          string
	PartyA             string
	IdentifierType     string
	Remarks            string
	QueueTimeOutURL    string
	ResultURL          string
}

// C2BRegisterURL is a model
type C2BRegisterURL struct {
	ShortCode       string
	ResponseType    string
	ConfirmationURL string
	ValidationURL   string
}

// ValidateRegisterURL validates the rquest register url fields.
func (m C2BRegisterURL) ValidateRegisterURL() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ResponseType, validation.Required),
		validation.Field(&m.ConfirmationURL, validation.Required),
		validation.Field(&m.ValidationURL, validation.Required),
	)
}

// C2BResponse ...
type C2BResponse struct {
	ID                int32 `json:"id,omitempty"`
	TransactionType   string
	TransID           string
	TransTime         string
	TransAmount       string
	BusinessShortCode string
	BillRefNumber     string
	InvoiceNumber     string
	OrgAccountBalance string
	ThirdPartyTransID string
	MSISDN            string
	FirstName         string
	MiddleName        string
	LastName          string
	Status            string
	ResultDesc        sql.NullString
	MeterNumber       string `json:"meter_number,omitempty"`
	MeterType         int `json:"meter_type,omitempty"`
	CompanyID         int    `json:"company_id"`
	AddedBy           int    `json:"added_by"`
}

// VResponse ...
type VResponse struct {
	ResultCode int
	ResultDesc string
}
