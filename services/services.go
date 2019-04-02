package services

import (
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// NewOwner ...
func NewOwner(m models.DeviceDetails) *models.VehicleOwner {
	vm := &models.VehicleOwner{
		OwnerID:    app.GenerateNewID(),
		OwnerIDNo:  m.OwnerID,
		OwnerName:  m.OwnerName,
		OwnerEmail: m.OwnerEmail,
		OwnerPhone: m.OwnerPhoneNumber,
	}

	return vm
}

// NewFitter
func NewFitter(m models.DeviceDetails) *models.FitterDetails {
	fd := &models.FitterDetails{
		FitterID:          app.GenerateNewID(),
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

// NewVehicle
func NewVehicle(m models.DeviceDetails) *models.VehicleDetails {
	vd := &models.VehicleDetails{
		VehicleID:       app.GenerateNewID(),
		VehicleStringID: strings.ToLower(strings.Replace(m.RegistrationNO, " ", "", -1)),
		VehicleRegNo:    m.RegistrationNO,
		ChassisNo:       m.ChasisNO,
		MakeType:        m.MakeType,
	}
	return vd
}

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
