package services

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// deviceDAO specifies the interface of the device DAO needed by DeviceService.
type deviceDAO interface {
	// Get returns the device with the specified device ID.
	Get(rs app.RequestScope, id int32) (*models.Devices, error)
	// Count returns the number of devices.
	Count(rs app.RequestScope, cid int) (int, error)
	// Query returns the list of devices with the given offset and limit.
	Query(rs app.RequestScope, offset, limit, cid int) ([]models.Devices, error)
	// QueryPositions returns the list of devices with their positions
	QueryPositions(rs app.RequestScope, offset, limit int, uid uint32, start, stop int64) ([]models.Devices, error)
	CountQueryPositions(rs app.RequestScope, uid uint32) (int, error)
	// Create saves a new device in the storage.
	Create(rs app.RequestScope, device *models.Devices) error
	// Update updates the device with given ID in the storage.
	Update(rs app.RequestScope, id int32, device *models.Devices) error
	// Delete removes the device with given ID from the storage.
	Delete(rs app.RequestScope, id int32) error
	CountConfiguredDevices(rs app.RequestScope, vehicleid, deviceid int) (int, error)
	ConfiguredDevices(rs app.RequestScope, offset, limit, vehicleid, deviceid int) ([]models.DeviceConfiguration, error)
}

// DeviceService provides services related with devices.
type DeviceService struct {
	dao deviceDAO
}

// NewDeviceService creates a new DeviceService with the given device DAO.
func NewDeviceService(dao deviceDAO) *DeviceService {
	return &DeviceService{dao}
}

// Get returns the device with the specified the device ID.
func (s *DeviceService) Get(rs app.RequestScope, id int32) (*models.Devices, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new device.
func (s *DeviceService) Create(rs app.RequestScope, model *models.Devices) (*models.Devices, error) {
	if err := model.ValidateDevices(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.ID)
}

// Update updates the device with the specified ID.
func (s *DeviceService) Update(rs app.RequestScope, id int32, model *models.Devices) (*models.Devices, error) {
	if err := model.ValidateDevices(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the device with the specified ID.
func (s *DeviceService) Delete(rs app.RequestScope, id int32) (*models.Devices, error) {
	device, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return device, err
}

// Count returns the number of devices.
func (s *DeviceService) Count(rs app.RequestScope, cid int) (int, error) {
	return s.dao.Count(rs, cid)
}

// CountQueryPositions returns the number of devices positions.
func (s *DeviceService) CountQueryPositions(rs app.RequestScope, uid uint32) (int, error) {
	return s.dao.CountQueryPositions(rs, uid)
}

// Query returns the devices with the specified offset and limit.
func (s *DeviceService) Query(rs app.RequestScope, offset, limit, cid int) ([]models.Devices, error) {
	return s.dao.Query(rs, offset, limit, cid)
}

// QueryPositions returns the device positions with the specified offset and limit.
func (s *DeviceService) QueryPositions(rs app.RequestScope, offset, limit int, uid uint32, start, stop int64) ([]models.Devices, error) {
	return s.dao.QueryPositions(rs, offset, limit, uid, start, stop)
}

// CountConfiguredDevices ...
func (s *DeviceService) CountConfiguredDevices(rs app.RequestScope, vehicleid, deviceid int) (int, error) {
	return s.dao.CountConfiguredDevices(rs, vehicleid, deviceid)
}

// ConfiguredDevices ...
func (s *DeviceService) ConfiguredDevices(rs app.RequestScope, offset, limit, vehicleid, deviceid int) ([]models.DeviceConfiguration, error) {
	return s.dao.ConfiguredDevices(rs, offset, limit, vehicleid, deviceid)
}
