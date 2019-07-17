package services

import (
	"strings"

	"github.com/ekas-portal-api/models"
)

// NewOwner ...
func NewOwner(m models.DeviceDetails, id uint32, userid uint32) *models.VehicleOwner {
	vm := &models.VehicleOwner{
		OwnerID:    id,
		UserID:     userid,
		OwnerIDNo:  m.OwnerID,
		OwnerName:  m.OwnerName,
		OwnerEmail: m.OwnerEmail,
		OwnerPhone: m.OwnerPhoneNumber,
	}

	return vm
}

// NewFitter ...
func NewFitter(m models.DeviceDetails, id uint32, userid uint32) *models.FitterDetails {
	fd := &models.FitterDetails{
		FitterID:          id,
		UserID:            userid,
		FitterIDNo:        m.AgentID,
		FittingCenterName: m.FittingCenter,
		FitterLocation:    m.AgentLocation,
		FitterEmail:       m.EmailAddress,
		FitterPhone:       m.AgentPhone,
		FittingDate:       m.FittingDate,
		FittingTime:       m.FittingTime,
		FitterBizRegNo:    m.BusinessRegNo,
	}

	return fd
}

// NewVehicle ...
func NewVehicle(m models.DeviceDetails, id uint32, userid uint32) *models.VehicleDetails {
	vd := &models.VehicleDetails{
		VehicleID:         id,
		UserID:            userid,
		VehicleStringID:   strings.ToLower(strings.Replace(m.RegistrationNO, " ", "", -1)),
		VehicleRegNo:      strings.ToUpper(m.RegistrationNO),
		ChassisNo:         strings.ToUpper(m.ChasisNO),
		MakeType:          strings.ToUpper(m.MakeType),
		NotificationEmail: strings.ToLower(m.NotificationEmail),
		NotificationNO:    m.NotificationNO,
	}
	return vd
}

// NewTrackingServerVehicle ...
func NewTrackingServerVehicle(m *models.Vehicle) *models.AddDeviceDetails {
	vd := &models.AddDeviceDetails{
		Name:               m.DeviceDetails.RegistrationNO,
		Imei:               m.GovernorDetails.DeviceID,
		IconID:             "45",
		FuelMeasurementID:  "1",
		TailLength:         "5",
		MinMovingSpeed:     "6",
		MinFuelFillings:    "10",
		MinFuelThefts:      "10",
		PlateNumber:        m.DeviceDetails.RegistrationNO,
		Vin:                m.DeviceDetails.ChasisNO,
		DeviceModel:        m.DeviceDetails.DeviceType,
		RegistrationNumber: "",
		ObjectOwner:        "",
	}
	return vd

}
