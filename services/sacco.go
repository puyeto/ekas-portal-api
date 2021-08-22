package services

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// saccoDAO specifies the interface of the sacco DAO needed by SaccoService.
type saccoDAO interface {
	// Get returns the sacco with the specified sacco ID.
	Get(rs app.RequestScope, id int) (*models.Saccos, error)
	// Get sacco details associated with a user
	GetSaccoUser(rs app.RequestScope, userid int) (*models.Saccos, error)
	// Count returns the number of saccos.
	Count(rs app.RequestScope) (int, error)
	// Query returns the list of saccos with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]models.Saccos, error)
	// Create saves a new sacco in the storage.
	Create(rs app.RequestScope, sacco *models.Saccos) error
	CreateSaccoUser(rs app.RequestScope, saccoid int32, userid int32) error
	// Update updates the sacco with given ID in the storage.
	Update(rs app.RequestScope, id int, sacco *models.Saccos) error
	// Delete removes the sacco with given ID from the storage.
	Delete(rs app.RequestScope, id int) error
	IsExistSaccoName(rs app.RequestScope, sacconame string) (int, error)
}

// SaccoService provides services related with saccos.
type SaccoService struct {
	dao saccoDAO
}

// NewSaccoService creates a new SaccoService with the given sacco DAO.
func NewSaccoService(dao saccoDAO) *SaccoService {
	return &SaccoService{dao}
}

// Get returns the sacco with the specified the sacco ID.
func (s *SaccoService) Get(rs app.RequestScope, id int) (*models.Saccos, error) {
	return s.dao.Get(rs, id)
}

// GetSaccoUser Get sacco details associated with a user
func (s *SaccoService) GetSaccoUser(rs app.RequestScope, userid int) (*models.Saccos, error) {
	return s.dao.GetSaccoUser(rs, userid)
}

// Create creates a new sacco.
func (s *SaccoService) Create(rs app.RequestScope, model *models.Saccos) (*models.Saccos, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}

	saccoid, err := s.dao.IsExistSaccoName(rs, model.Name)
	if err != nil {
		return nil, err
	}

	if saccoid > 0 {
		if err := s.dao.Update(rs, saccoid, model); err != nil {
			return nil, err
		}
	} else {
		if err := s.dao.Create(rs, model); err != nil {
			return nil, err
		}
	}

	return s.dao.Get(rs, int(model.ID))
}

// Update updates the sacco with the specified ID.
func (s *SaccoService) Update(rs app.RequestScope, id int, model *models.Saccos) (*models.Saccos, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}

	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the sacco with the specified ID.
func (s *SaccoService) Delete(rs app.RequestScope, id int) (*models.Saccos, error) {
	sacco, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return sacco, err
}

// Count returns the number of saccos.
func (s *SaccoService) Count(rs app.RequestScope) (int, error) {
	return s.dao.Count(rs)
}

// Query returns the saccos with the specified offset and limit.
func (s *SaccoService) Query(rs app.RequestScope, offset, limit int) ([]models.Saccos, error) {
	return s.dao.Query(rs, offset, limit)
}
