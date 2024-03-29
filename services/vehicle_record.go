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
	Count(rs app.RequestScope, uid int, typ string, userdetails models.AuthUsers) (int, error)
	CountFilter(rs app.RequestScope, m *models.FilterVehicles) (int, error)
	// Query returns the list of vehicleRecords with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int, uid int, typ string, userdetails models.AuthUsers) ([]models.VehicleDetails, error)
	QueryFilter(rs app.RequestScope, offset, limit int, m *models.FilterVehicles) ([]models.VehicleDetails, error)
	// Update updates the vehicleRecord with given ID in the storage.
	UpdateVehicle(rs app.RequestScope, vehicleRecord *models.VehicleDetails) error
	// Delete removes the vehicleRecord with given ID from the storage.
	Delete(rs app.RequestScope, id uint32) error
	// CreateVehicle create new vehicle
	CreateVehicle(rs app.RequestScope, model *models.VehicleDetails) (uint32, error)
	ListVehicleRenewals(rs app.RequestScope, offset, limit int) ([]models.VehicleRenewals, error)
	RenewVehicle(rs app.RequestScope, model *models.VehicleRenewals) (uint32, error)
	CountRenewals(rs app.RequestScope) (int, error)
	CreateReminder(rs app.RequestScope, model *models.Reminders) (uint32, error)
	CountReminders(rs app.RequestScope, uid int) (int, error)
	GetReminder(rs app.RequestScope, offset, limit int, uid int) ([]models.Reminders, error)
	// GetUser returns the user with the specified user ID.
	GetUser(rs app.RequestScope, id uint32) (models.AuthUsers, error)
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
func (s *VehicleRecordService) Update(rs app.RequestScope, model *models.VehicleDetails) error {
	model.VehicleStringID = strings.ToLower(strings.Replace(model.VehicleRegNo, " ", "", -1))
	if model.Manufacturer == "" {
		model.Manufacturer = model.MakeType
	}

	model.Prepare("update")
	if err := model.ValidateVehicleDetails(); err != nil {
		return err
	}

	if err := s.dao.UpdateVehicle(rs, model); err != nil {
		return err
	}
	return nil
}

// Delete deletes the vehicleRecord with the specified ID.
func (s *VehicleRecordService) Delete(rs app.RequestScope, id uint32) error {
	return s.dao.Delete(rs, id)
}

// Count returns the number of vehicleRecords.
func (s *VehicleRecordService) Count(rs app.RequestScope, uid int, typ string, userdetails models.AuthUsers) (int, error) {
	return s.dao.Count(rs, uid, typ, userdetails)
}

// CountFilter returns the number of filtered vehicleRecords.
func (s *VehicleRecordService) CountFilter(rs app.RequestScope, model *models.FilterVehicles) (int, error) {
	return s.dao.CountFilter(rs, model)
}

// Query returns the vehicleRecords with the specified offset and limit.
func (s *VehicleRecordService) Query(rs app.RequestScope, offset, limit int, uid int, typ string, userdetails models.AuthUsers) ([]models.VehicleDetails, error) {
	return s.dao.Query(rs, offset, limit, uid, typ, userdetails)
}

// QueryFilter returns the filtered vehicleRecords with the specified offset and limit.
func (s *VehicleRecordService) QueryFilter(rs app.RequestScope, offset, limit int, model *models.FilterVehicles) ([]models.VehicleDetails, error) {
	return s.dao.QueryFilter(rs, offset, limit, model)
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

// RenewVehicle renew a vehicle.
func (s *VehicleRecordService) RenewVehicle(rs app.RequestScope, model *models.VehicleRenewals) (uint32, error) {
	if err := model.Validate(); err != nil {
		return 0, err
	}
	return s.dao.RenewVehicle(rs, model)
}

// ListVehicleRenewals renew a vehicle.
func (s *VehicleRecordService) ListVehicleRenewals(rs app.RequestScope, offset, limit int) ([]models.VehicleRenewals, error) {
	return s.dao.ListVehicleRenewals(rs, offset, limit)
}

// CreateReminder creates a new reminder.
func (s *VehicleRecordService) CreateReminder(rs app.RequestScope, model *models.Reminders) (uint32, error) {
	if err := model.ValidateReminders(); err != nil {
		return 0, err
	}

	return s.dao.CreateReminder(rs, model)
}

// CountRenewals returns the number of vehicleRecords.
func (s *VehicleRecordService) CountRenewals(rs app.RequestScope) (int, error) {
	return s.dao.CountRenewals(rs)
}

// CountReminders returns the number of reminderRecords.
func (s *VehicleRecordService) CountReminders(rs app.RequestScope, uid int) (int, error) {
	return s.dao.CountReminders(rs, uid)
}

// GetReminder returns the reminderRecords with the specified offset and limit.
func (s *VehicleRecordService) GetReminder(rs app.RequestScope, offset, limit int, uid int) ([]models.Reminders, error) {
	return s.dao.GetReminder(rs, offset, limit, uid)
}

// GetUser returns the user with the specified the user ID.
func (u *VehicleRecordService) GetUser(rs app.RequestScope, id uint32) (models.AuthUsers, error) {
	return u.dao.GetUser(rs, id)
}
