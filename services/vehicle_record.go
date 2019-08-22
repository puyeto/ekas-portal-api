package services

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// vehicleRecordDAO specifies the interface of the vehicleRecord DAO needed by VehicleRecordService.
type vehicleRecordDAO interface {
	// Get returns the vehicleRecord with the specified vehicleRecord ID.
	Get(rs app.RequestScope, id uint32) (*models.VehicleDetails, error)
	// Count returns the number of vehicleRecords.
	Count(rs app.RequestScope, uid int) (int, error)
	// Query returns the list of vehicleRecords with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int, uid int) ([]models.VehicleDetails, error)
	// Update updates the vehicleRecord with given ID in the storage.
	Update(rs app.RequestScope, id uint32, vehicleRecord *models.VehicleDetails) error
	// Delete removes the vehicleRecord with given ID from the storage.
	Delete(rs app.RequestScope, id uint32) error
	// CreateVehicle create new vehicle
	CreateVehicle(rs app.RequestScope, model *models.VehicleDetails) (uint32, error)
}

// VehicleRecordService provides services related with vehicleRecords.
type VehicleRecordService struct {
	dao vehicleRecordDAO
}

// NewVehicleRecordService creates a new VehicleRecordService with the given vehicleRecord DAO.
func NewVehicleRecordService(dao vehicleRecordDAO) *VehicleRecordService {
	return &VehicleRecordService{dao}
}

// Get returns the vehicleRecord with the specified the vehicleRecord ID.
func (s *VehicleRecordService) Get(rs app.RequestScope, id uint32) (*models.VehicleDetails, error) {
	return s.dao.Get(rs, id)
}

// Update updates the vehicleRecord with the specified ID.
func (s *VehicleRecordService) Update(rs app.RequestScope, id uint32, model *models.VehicleDetails) (*models.VehicleDetails, error) {
	if err := model.ValidateVehicleDetails(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the vehicleRecord with the specified ID.
func (s *VehicleRecordService) Delete(rs app.RequestScope, id uint32) (*models.VehicleDetails, error) {
	vehicleRecord, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return vehicleRecord, err
}

// Count returns the number of vehicleRecords.
func (s *VehicleRecordService) Count(rs app.RequestScope, uid int) (int, error) {
	return s.dao.Count(rs, uid)
}

// Query returns the vehicleRecords with the specified offset and limit.
func (s *VehicleRecordService) Query(rs app.RequestScope, offset, limit int, uid int) ([]models.VehicleDetails, error) {
	return s.dao.Query(rs, offset, limit, uid)
}

// CreateVehicle creates a new vehicle.
func (s *VehicleRecordService) CreateVehicle(rs app.RequestScope, model *models.VehicleDetails) (uint32, error) {
	if err := model.ValidateVehicleDetails(); err != nil {
		return 0, err
	}
	return s.dao.CreateVehicle(rs, model)
}
