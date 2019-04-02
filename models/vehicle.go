package models

// Vehicle ...
type Vehicle struct {
	VehicleDetails  TrackingVehicleDetails `json:"dehicle_details,omitempty"`
	DeviceDetails   DeviceDetails          `json:"device_detail"`
	GovernorDetails GovernorDetails        `json:"governor_details"`
	UserID          string                 `json:"user_id,omitempty"`
	UserHash        string                 `json:"user_hash,omitempty"`
	SimNO           string                 `json:"sim_no,omitempty"`
	SimIMEI         string                 `json:"sim_imei,omitempty"`
}

// TrackingVehicleDetails ...
type TrackingVehicleDetails struct{}

// DeviceDetails ....
type DeviceDetails struct {
	OwnerName        string `json:"owner_name"`
	OwnerID          string `json:"owner_id"`
	OwnerPhoneNumber string `json:"owner_phone_number,omitempty"`
	OwnerEmail       string `json:"owner_email,omitempty"`
	RegistrationNO   string `json:"registration_no"`
	ChasisNO         string `json:"chasis_no"`
	MakeType         string `json:"make_type"`
	Certificate      string `json:"certificate"`
	DeviceType       string `json:"device_type"`
	SerialNO         string `json:"serial_no"`
	FittingDate      string `json:"fitting_date"`
	FittingTime      string `json:"fitting_time"`
	FittingCenter    string `json:"fitting_center"`
	AgentID          string `json:"agent_id"`
	AgentLocation    string `json:"agent_location"`
	EmailAddress     string `json:"email_address"`
	AgentPhone       string `json:"agent_phone"`
	BusinessRegNo    string `json:"business_reg_no"`
	SetAlarm         string `json:"set_alarm"`
	SetFrequency     string `json:"set_frequency"`
	PresetSpeed      string `json:"preset_speed"`
	GPRSSetSpeed     string `json:"gprs_set_speed"`
	SpeedSource      string `json:"speed_source"`
	ConfigDone       string `json:"config_done"`
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
	VehicleID       uint64 `json:"vehicle_id" db:"pk,vehicle_id"`
	VehicleStringID string `json:"vehicle_string_id,omitempty" db:"vehicle_string_id"`
	VehicleRegNo    string `json:"vehicle_reg_no" db:"vehicle_reg_no"`
	ChassisNo       string `json:"chassis_no" db:"chassis_no"`
	MakeType        string `json:"make_type" db:"make_type"`
}

// VehicleOwner ...
type VehicleOwner struct {
	OwnerID    uint64 `json:"owner_id" db:"pk,owner_id"`
	OwnerIDNo  string `json:"owner_id_no" db:"owner_id_no"`
	OwnerName  string `json:"owner_name" db:"owner_name"`
	OwnerEmail string `json:"owner_email" db:"owner_email"`
	OwnerPhone string `json:"owner_phone" db:"owner_phone"`
}

// FitterDetails ...
type FitterDetails struct {
	FitterID          uint64 `json:"fitting_id" db:"pk,fitting_id"`
	FitterIDNo        string `json:"fitting_id_no" db:"fitting_id_no"`
	FittingCenterName string `json:"fitting_center_name" db:"fitting_center_name"`
	FitterLocation    string `json:"fitter_location" db:"fitter_location"`
	FitterEmail       string `json:"fitter_email,omitempty" db:"fitter_email"`
	FitterAddress     string `json:"fitter_address,omitempty" db:"fitter_address"`
	FitterPhone       string `json:"fitter_phone" db:"fitter_phone"`
	FittingDate       string `json:"fitting_date" db:"fitting_date"`
	FittingTime       string `json:"fitting_time" db:"fitting_time"`
	FitterBizRegNo    string `json:"fitter_biz_reg_no" db:"fitter_biz_reg_no"`
}
