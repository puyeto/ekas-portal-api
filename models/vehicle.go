package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Vehicle ...
type Vehicle struct {
	DeviceDetails   DeviceDetails   `json:"device_detail"`
	GovernorDetails GovernorDetails `json:"governor_details"`
	UserID          uint32          `json:"user_id,omitempty"`
	UserHash        string          `json:"user_hash,omitempty"`
	SimNO           string          `json:"sim_no,omitempty"`
	SimIMEI         string          `json:"sim_imei,omitempty"`
	VehicleID       uint32          `json:"vehicle_id,omitempty"`
	OwnerID         uint32          `json:"owner_id,omitempty"`
	FitterID        uint32          `json:"fitter_id,omitempty"`
	ConfigID        uint32          `json:"conf_id,omitempty"`
}

// VehicleConfigDetails ...
type VehicleConfigDetails struct {
	ConfigID          uint32 `json:"conf_id"`
	DeviceID          uint32 `json:"device_id,omitempty"`
	VehicleID         uint32 `json:"vehicle_id,omitempty"`
	OwnerID           uint32 `json:"owner_id,omitempty"`
	FitterID          uint32 `json:"fitter_id,omitempty"`
	NotificationEmail string `json:"notification_email,omitempty"`
	NotificationNO    string `json:"notification_no,omitempty"`
	Data              string `json:"vehicle_data,omitempty"`
	CreatedOn         string `json:"created_on,omitempty"`
}

// SearchDetails ...
type SearchDetails struct {
	VehicleName string `json:"vehicle_name"`
	Data        string `json:"vehicle_data,omitempty"`
}

// DeviceDetails ....
type DeviceDetails struct {
	OwnerName         string `json:"owner_name"`
	OwnerID           string `json:"owner_id"`
	OwnerPhoneNumber  string `json:"owner_phone_number,omitempty"`
	OwnerEmail        string `json:"owner_email,omitempty"`
	RegistrationNO    string `json:"registration_no"`
	ChasisNO          string `json:"chasis_no"`
	MakeType          string `json:"make_type"`
	Certificate       string `json:"certificate"`
	DeviceType        string `json:"device_type"`
	SerialNO          string `json:"serial_no"`
	FittingDate       string `json:"fitting_date"`
	FittingTime       string `json:"fitting_time"`
	FittingCenter     string `json:"fitting_center"`
	AgentID           int    `json:"agent_id"`
	AgentLocation     string `json:"agent_location"`
	EmailAddress      string `json:"email_address"`
	AgentPhone        string `json:"agent_phone"`
	BusinessRegNo     string `json:"business_reg_no"`
	SetAlarm          string `json:"set_alarm"`
	SetFrequency      string `json:"set_frequency"`
	PresetSpeed       string `json:"preset_speed"`
	GPRSSetSpeed      string `json:"gprs_set_speed"`
	SpeedSource       string `json:"speed_source"`
	ConfigDone        string `json:"config_done"`
	NotificationEmail string `json:"notification_email"`
	NotificationNO    string `json:"notification_no"`
}

// GovernorDetails ...
type GovernorDetails struct {
	DeviceID       string `json:"device_id"`
	AccountID      string `json:"account_id"`
	Domain         string `json:"domain"`
	Port           string `json:"port"`
	SecondDomain   string `json:"second_domain"`
	SecondPort     string `json:"second_port"`
	FailSafe       string `json:"fail_safe"`
	APN            string `json:"apn"`
	APNSet         string `json:"apn_set"`
	APNUsername    string `json:"apn_username"`
	APNUsernamSet  string `json:"apn_username_set"`
	APNPassword    string `json:"apn_password"`
	APNPasswordSet string `json:"apn_password_set"`
	TConfigDone    string `json:"t_config_done"`
}

// VehicleDetails ...
type VehicleDetails struct {
	VehicleID              uint32    `json:"vehicle_id" db:"pk,vehicle_id"`
	UserID                 uint32    `json:"user_id" db:"user_id"`
	OwnerID                uint32    `json:"owner_id" db:"owner_id"`
	CompanyID              uint32    `json:"company_id" db:"company_id"`
	DeviceID              uint32    `json:"device_id" db:"device_id"`
	CompanyName            string    `json:"company_name,omitempty"`
	VehicleStringID        string    `json:"vehicle_string_id,omitempty" db:"vehicle_string_id"`
	VehicleRegNo           string    `json:"vehicle_reg_no" db:"vehicle_reg_no"`
	ChassisNo              string    `json:"chassis_no" db:"chassis_no"`
	MakeType               string    `json:"make_type" db:"make_type"`
	NotificationEmail      string    `json:"notification_email,omitempty" db:"notification_email"`
	NotificationNO         string    `json:"notification_no,omitempty" db:"notification_no"`
	VehicleStatus          int8      `json:"status" db:"vehicle_status"`
	AutoInvoicing          int8      `json:"auto_invoicing,omitempty" db:"auto_invoicing"`
	InvoiceDueDate         time.Time `json:"invoice_due_date,omitempty" db:"invoice_due_date"`
	CreatedOn              time.Time `json:"created_on" db:"created_on"`
	Model                  string    `json:"model,omitempty" db:"model"`
	ModelYear              int16     `json:"model_year,omitempty" db:"model_year"`
	Manufacturer           string    `json:"manufacturer,omitempty" db:"manufacturer"`
	BodyStyle              string    `json:"body_style,omitempty" db:"body_style"`
	BodyType               string    `json:"body_type,omitempty" db:"body_type"`
	DeleteTripDetailsAfter string    `json:"delete_trip_details_after,omitempty" db:"delete_trip_details_after"`
	DeleteTripsAfter       string    `json:"delete_trips_after,omitempty" db:"delete_trips_after"`
	FuelType               int       `json:"fuel_type,omitempty" db:"fuel_type"`
	DefaultTripType        int       `json:"default_trip_type,omitempty" db:"default_trip_type"`
}

// ValidateVehicleDetails validates fields.
func (v VehicleDetails) ValidateVehicleDetails() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.VehicleRegNo, validation.Required),
		validation.Field(&v.UserID, validation.Required),
		validation.Field(&v.ChassisNo, validation.Required),
		validation.Field(&v.MakeType, validation.Required),
	)
}

// FitterDetails ...
type FitterDetails struct {
	FitterID          uint32 `json:"fitting_id" db:"pk,fitting_id"`
	UserID            uint32 `json:"user_id" db:"user_id"`
	FitterIDNo        int    `json:"fitting_id_no" db:"fitting_id_no"`
	FittingCenterName string `json:"fitting_center_name" db:"fitting_center_name"`
	FitterLocation    string `json:"fitter_location" db:"fitter_location"`
	FitterEmail       string `json:"fitter_email,omitempty" db:"fitter_email"`
	FitterAddress     string `json:"fitter_address,omitempty" db:"fitter_address"`
	FitterPhone       string `json:"fitter_phone" db:"fitter_phone"`
	FittingDate       string `json:"fitting_date" db:"fitting_date"`
	FittingTime       string `json:"fitting_time" db:"fitting_time"`
	FitterBizRegNo    string `json:"fitter_biz_reg_no" db:"fitter_biz_reg_no"`
}

// TripData ...
type TripData struct {
	TripID     int32     `json:"trip_id" db:"pk,trip_id"`
	DeviceID   int32     `json:"device_id" db:"device_id"`
	DateTime   time.Time `json:"date_time" db:"data_date"`
	Speed      float32   `json:"speed,omitempty" db:"speed"`
	Longitude  float32   `json:"longitude,omitempty" db:"longitude"`
	Latitude   float32   `json:"latitude,omitempty" db:"latitude"`
	Failsafe   bool      `json:"failsafe,omitempty" db:"failsafe"`
	Disconnect bool      `json:"disconnect,omitempty" db:"disconnect"`
}

// TripBetweenDates ...
type TripBetweenDates struct {
	DeviceID string `json:"DeviceID,omitempty"`
	From     int64  `json:"From,omitempty"`
	To       int64  `json:"To,omitempty"`
}

// DeviceData ...
type DeviceData struct {
	SystemCode                     string    `json:"system_code,omitempty"`                      // 4 bytes
	SystemMessage                  int       `json:"system_message,omitempty"`                   // 1 byte
	DeviceID                       uint32    `json:"device_id,omitempty"`                        // 4 bytes
	CommunicationControlField      uint32    `json:"communication_control_field,omitempty"`      // 2 bytes
	MessageNumerator               int       `json:"message_numerator,omitempty"`                // 1 byte
	HardwareVersion                int       `json:"hardware_version,omitempty"`                 // 1 byte
	SoftwareVersion                int       `json:"software_version,omitempty"`                 // 1 byte
	ProtocolVersionIdentifier      int       `json:"protocol_version_identifier,omitempty"`      // 1 byte
	Status                         int       `json:"status,omitempty"`                           // 1 byte
	ConfigurationFlags             int       `json:"configuration_flags,omitempty"`              // 2 bytes
	TransmissionReasonSpecificData int       `json:"transmission_reason_specificData,omitempty"` // 1 byte
	Failsafe                       bool      `json:"failsafe"`
	Disconnect                     bool      `json:"disconnect"`
	Offline                        bool      `json:"offline"`
	TransmissionReason             int       `json:"transmission_reason,omitempty"` // 1 byte
	ModeOfOperation                int       `json:"mode_of_operation,omitempty"`   // 1 byte
	IOStatus                       uint16    `json:"io_status,omitempty"`           // 5 bytes
	AnalogInput1Value              uint16    `json:"analog_Input_1_value,omitempty"`
	AnalogInput1Value1             uint16    `json:"analog_Input_1_value_1,omitempty"`
	AnalogInput2Value              uint16    `json:"analog_Input_2_value,omitempty"`
	AnalogInput2Value2             uint16    `json:"analog_Input_2_value_2,omitempty"`
	MileageCounter                 uint16    `json:"mileage_counter,omitempty"` // 3 bytes
	DriverID                       uint16    `json:"driver_id,omitempty"`       // 6 bytes
	LastGPSFix                     uint16    `json:"last_gps_fix,omitempty"`
	LocationStatus                 uint16    `json:"location_status,omitempty"`
	Mode1                          uint16    `json:"mode_1,omitempty"`
	Mode2                          uint16    `json:"mode_2,omitempty"`
	NoOfSatellitesUsed             int       `json:"no_of_satellites_used,omitempty"` // 1 byte
	Longitude                      int32     `json:"longitude,omitempty"`             // 4 byte
	Latitude                       int32     `json:"latitude,omitempty"`              // 4 byte
	Altitude                       int32     `json:"altitude,omitempty"`              // 4 byte
	GroundSpeed                    float32    `json:"ground_speed,omitempty"`          // 4 byte
	SpeedDirection                 int       `json:"speed_direction,omitempty"`       // 2 byte
	UTCTimeSeconds                 int       `json:"utc_time_seconds,omitempty"`      // 1 byte
	UTCTimeMinutes                 int       `json:"utc_time_minutes,omitempty"`      // 1 byte
	UTCTimeHours                   int       `json:"utc_time_hours,omitempty"`        // 1 byte
	UTCTimeDay                     int       `json:"utc_time_day,omitempty"`          // 1 byte
	UTCTimeMonth                   int       `json:"utc_time_month,omitempty"`        // 1 byte
	UTCTimeYear                    int       `json:"utc_time_year,omitempty"`         // 2 byte
	ErrorDetectionCode             uint16    `json:"error_detection_code,omitempty"`
	DateTime                       time.Time `json:"date_time,omitempty"`
	Name                           string    `json:"name,omitempty"`
	DateTimeStamp                  int64     `json:"date_time_stamp,omitempty"`
}

// LastSeenStruct ...
type LastSeenStruct struct {
	DateTime   time.Time
	DeviceData DeviceData
}

// CurrentViolations ...
type CurrentViolations struct {
	DeviceID            string    `json:"device_id" db:"device_id"`
	Name                string    `json:"name" db:"name"`
	UserID              string    `json:"user_id,omitempty" db:"user_id"`
	OverspeedTripData   string    `json:"overspeed_trip_data,omitempty" db:"overspeed_trip_data"`
	OverspeedSpeed      string    `json:"overspeed_speed,omitempty" db:"overspeed_speed"`
	OverspeedDate       string    `json:"overspeed_date,omitempty" db:"overspeed_date"`
	DisconnectTripData  string    `json:"disconnect_trip_data,omitempty" db:"disconnect_trip_data"`
	DisconnectTripSpeed string    `json:"disconnect_trip_speed,omitempty" db:"disconnect_trip_speed"`
	DisconnectTripDate  string    `json:"disconnect_trip_date,omitempty" db:"disconnect_trip_date"`
	FailsafeTripData    string    `json:"failsafe_trip_data,omitempty" db:"failsafe_trip_data"`
	FailsafeTripSpeed   string    `json:"failsafe_trip_speed,omitempty" db:"failsafe_trip_speed"`
	FailsafeTripDate    string    `json:"failsafe_trip_date,omitempty" db:"failsafe_trip_date"`
	OfflineTripData     string    `json:"offline_trip_data,omitempty" db:"offline_trip_data"`
	OfflineTripSpeeds   string    `json:"offline_trip_speed,omitempty" db:"offline_trip_speed"`
	OfflineTripDate     string    `json:"offline_trip_date,omitempty" db:"offline_trip_date"`
	CreatedOn           time.Time `json:"created_on,omitempty" db:"created_on"`
}

// Reminders ...
type Reminders struct {
	ID     uint32 `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	UserID uint32 `json:"user_id" db:"user_id"`
}

// ValidateReminders validates fields.
func (r Reminders) ValidateReminders() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
	)
}
