package services

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// ownerDAO specifies the interface of the owner DAO needed by OwnerService.
type ownerDAO interface {
	// Get returns the owner with the specified owner ID.
	Get(rs app.RequestScope, id uint32) (*models.VehicleOwner, error)
	// Count returns the number of owners.
	Count(rs app.RequestScope) (int, error)
	// Query returns the list of owners with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]models.VehicleOwner, error)
	// Create saves a new owner in the storage.
	Create(rs app.RequestScope, owner *models.VehicleOwner) error
	// Update updates the owner with given ID in the storage.
	Update(rs app.RequestScope, id uint32, owner *models.VehicleOwner) error
	// Delete removes the owner with given ID from the storage.
	Delete(rs app.RequestScope, id uint32) error
}

// OwnerService provides services related with owners.
type OwnerService struct {
	dao ownerDAO
}

// NewOwnerService creates a new OwnerService with the given owner DAO.
func NewOwnerService(dao ownerDAO) *OwnerService {
	return &OwnerService{dao}
}

// Get returns the owner with the specified the owner ID.
func (s *OwnerService) Get(rs app.RequestScope, id uint32) (*models.VehicleOwner, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new owner.
func (s *OwnerService) Create(rs app.RequestScope, model *models.VehicleOwner) (*models.VehicleOwner, error) {
	if err := model.ValidateVehicleOwner(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.OwnerID)
}

// Update updates the owner with the specified ID.
func (s *OwnerService) Update(rs app.RequestScope, id uint32, model *models.VehicleOwner) (*models.VehicleOwner, error) {
	if err := model.ValidateVehicleOwner(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the owner with the specified ID.
func (s *OwnerService) Delete(rs app.RequestScope, id uint32) (*models.VehicleOwner, error) {
	owner, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return owner, err
}

// Count returns the number of owners.
func (s *OwnerService) Count(rs app.RequestScope) (int, error) {
	return s.dao.Count(rs)
}

// Query returns the owners with the specified offset and limit.
func (s *OwnerService) Query(rs app.RequestScope, offset, limit int) ([]models.VehicleOwner, error) {
	return s.dao.Query(rs, offset, limit)
}
