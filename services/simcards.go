package services

import (
	"errors"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// simcardDAO specifies the interface of the simcard DAO needed by SimcardService.
type simcardDAO interface {
	// Get returns the simcard with the specified simcard ID.
	Get(rs app.RequestScope, id int) (*models.Simcards, error)
	// Count returns the number of simcards.
	Count(rs app.RequestScope, status string) (int, error)
	// Query returns the list of simcards with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int, status string) ([]models.Simcards, error)
	// Create saves a new simcard in the storage.
	Create(rs app.RequestScope, simcard *models.Simcards) error
	// Update updates the simcard with given ID in the storage.
	Update(rs app.RequestScope, id int, simcard *models.Simcards) error
	// Delete removes the simcard with given ID from the storage.
	Delete(rs app.RequestScope, id int) error
	IsExistSimcard(rs app.RequestScope, simcardname string) (int, error)
	GetStats(rs app.RequestScope) *models.SimcardStats
}

// SimcardService provides services related with simcards.
type SimcardService struct {
	dao simcardDAO
}

// NewSimcardService creates a new SimcardService with the given simcard DAO.
func NewSimcardService(dao simcardDAO) *SimcardService {
	return &SimcardService{dao}
}

// Get returns the simcard with the specified the simcard ID.
func (s *SimcardService) Get(rs app.RequestScope, id int) (*models.Simcards, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new simcard.
func (s *SimcardService) Create(rs app.RequestScope, model *models.Simcards) (*models.Simcards, error) {
	model.Prepare()
	if err := model.Validate(); err != nil {
		return nil, err
	}

	simcardid, err := s.dao.IsExistSimcard(rs, model.Identifier)
	if err != nil {
		return nil, err
	}

	if simcardid > 0 {
		return nil, errors.New("Simcard already exist")
	}

	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}

	return s.dao.Get(rs, int(model.ID))
}

// Update updates the simcard with the specified ID.
func (s *SimcardService) Update(rs app.RequestScope, id int, model *models.Simcards) (*models.Simcards, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}

	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the simcard with the specified ID.
func (s *SimcardService) Delete(rs app.RequestScope, id int) (*models.Simcards, error) {
	simcard, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return simcard, err
}

// Count returns the number of simcards.
func (s *SimcardService) Count(rs app.RequestScope, status string) (int, error) {
	return s.dao.Count(rs, status)
}

// Query returns the simcards with the specified offset and limit.
func (s *SimcardService) Query(rs app.RequestScope, offset, limit int, status string) ([]models.Simcards, error) {
	return s.dao.Query(rs, offset, limit, status)
}

// GetStats ...
func (s *SimcardService) GetStats(rs app.RequestScope) *models.SimcardStats {
	return s.dao.GetStats(rs)
}
