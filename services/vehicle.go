package services

import (
	"fmt"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// vehicleDAO specifies the interface of the vehicle DAO needed by VehicleService.
type vehicleDAO interface {
	GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error)
	// Create saves a new vehicle in the storage.
	CreateVehicle(rs app.RequestScope, vehicle *models.VehicleDetails) error
	CreateVehicleOwner(rs app.RequestScope, vo *models.VehicleOwner) error
	CreateFitter(rs app.RequestScope, fd *models.FitterDetails) error
	CreateConfiguration(rs app.RequestScope, vehicle *models.Vehicle, ownerid uint32, fitterid uint32, vehicleid uint32) error
	UpdateConfigurationStatus(rs app.RequestScope, configid uint32, status int8) error
}

// VehicleService provides services related with vehicles.
type VehicleService struct {
	dao vehicleDAO
}

// NewVehicleService creates a new VehicleService with the given vehicle DAO.
func NewVehicleService(dao vehicleDAO) *VehicleService {
	return &VehicleService{dao}
}

// GetVehicleByStrID ...
func (s *VehicleService) GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error) {
	return s.dao.GetVehicleByStrID(rs, strid)
}

// Create creates a new vehicle.
func (s *VehicleService) Create(rs app.RequestScope, model *models.Vehicle) (int, error) {
	// if err := model.Validate(); err != nil {
	//	return nil, err
	// }

	// Add vehicle owner
	var ownerid = app.GenerateNewID()
	if model.OwnerID > 0 {
		ownerid = model.OwnerID
	}
	vm := NewOwner(model.DeviceDetails, ownerid)
	if err := s.dao.CreateVehicleOwner(rs, vm); err != nil {
		return 0, err
	}

	// Add Fitter Center / Fitter
	var fid = app.GenerateNewID()
	if model.FitterID > 0 {
		fid = model.FitterID
	}
	fd := NewFitter(model.DeviceDetails, fid)
	if err := s.dao.CreateFitter(rs, fd); err != nil {
		return 0, err
	}

	// Add Vehicle
	fmt.Println(model.VehicleID)
	var vid = app.GenerateNewID()
	if model.VehicleID > 0 {
		vid = model.VehicleID
	}
	fmt.Println(vid)
	vd := NewVehicle(model.DeviceDetails, vid)
	if err := s.dao.CreateVehicle(rs, vd); err != nil {
		return 0, err
	}

	// Add Configuartion Details
	fmt.Println(model.ConfigID)
	if model.ConfigID > 0 {
		// update configuration status
		if err := s.dao.UpdateConfigurationStatus(rs, model.ConfigID, 0); err != nil {
			return 0, err
		}
	}
	if err := s.dao.CreateConfiguration(rs, model, vm.OwnerID, fd.FitterID, vd.VehicleID); err != nil {
		return 0, err
	}

	// Add vehicle to tracking server
	tsv := NewTrackingServerVehicle(model)
	_, err := AddDevicesTrackingServer(rs, tsv, "en", model.UserHash)
	if err != nil {
		return 0, err
	}

	return 0, nil
}
