package services

import (
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// vehicleRecordDAO specifies the interface of the vehicleRecord DAO needed by VehicleRecordService.
type vehicleRecordDAO interface {
	// Get returns the vehicleRecord with the specified vehicleRecord ID.
	Get(rs app.RequestScope, id uint32) (*models.VehicleDetails, error)
	// Count returns the number of vehicleRecords.
	Count(rs app.RequestScope, uid int, typ string) (int, error)
	// Query returns the list of vehicleRecords with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int, uid int, typ string) ([]models.VehicleDetails, error)
	// Update updates the vehicleRecord with given ID in the storage.
	Update(rs app.RequestScope, id uint32, vehicleRecord *models.VehicleDetails) error
	// Delete removes the vehicleRecord with given ID from the storage.
	Delete(rs app.RequestScope, id uint32) error
	// CreateVehicle create new vehicle
	CreateVehicle(rs app.RequestScope, model *models.VehicleDetails) (uint32, error)
	CreateReminder(rs app.RequestScope, model *models.Reminders) (uint32, error)
	CountReminders(rs app.RequestScope, uid int) (int, error)
	GetReminder(rs app.RequestScope, offset, limit int, uid int) ([]models.Reminders, error)
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
func (s *VehicleRecordService) Delete(rs app.RequestScope, id uint32) error {
	return s.dao.Delete(rs, id)
}

// Count returns the number of vehicleRecords.
func (s *VehicleRecordService) Count(rs app.RequestScope, uid int, typ string) (int, error) {
	return s.dao.Count(rs, uid, typ)
}

// Query returns the vehicleRecords with the specified offset and limit.
func (s *VehicleRecordService) Query(rs app.RequestScope, offset, limit int, uid int, typ string) ([]models.VehicleDetails, error) {
	return s.dao.Query(rs, offset, limit, uid, typ)
}

// CreateVehicle creates a new vehicle.
func (s *VehicleRecordService) CreateVehicle(rs app.RequestScope, model *models.VehicleDetails) (uint32, error) {
	if err := model.ValidateVehicleDetails(); err != nil {
		return 0, err
	}
	model.VehicleStringID = strings.ToLower(strings.Replace(model.VehicleRegNo, " ", "", -1))
	if model.Manufacturer == "" {
		model.Manufacturer = model.MakeType
	}
	return s.dao.CreateVehicle(rs, model)
}

// CreateReminder creates a new reminder.
func (s *VehicleRecordService) CreateReminder(rs app.RequestScope, model *models.Reminders) (uint32, error) {
	if err := model.ValidateReminders(); err != nil {
		return 0, err
	}

	return s.dao.CreateReminder(rs, model)
}

// CountReminders returns the number of reminderRecords.
func (s *VehicleRecordService) CountReminders(rs app.RequestScope, uid int) (int, error) {
	return s.dao.CountReminders(rs, uid)
}

// GetReminder returns the reminderRecords with the specified offset and limit.
func (s *VehicleRecordService) GetReminder(rs app.RequestScope, offset, limit int, uid int) ([]models.Reminders, error) {
	return s.dao.GetReminder(rs, offset, limit, uid)
}
