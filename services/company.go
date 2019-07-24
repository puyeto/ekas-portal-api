package services

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// companyDAO specifies the interface of the company DAO needed by CompanyService.
type companyDAO interface {
	// Get returns the company with the specified company ID.
	Get(rs app.RequestScope, id int) (*models.Companies, error)
	// Count returns the number of companys.
	Count(rs app.RequestScope) (int, error)
	// Query returns the list of companys with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]models.Companies, error)
	// Create saves a new company in the storage.
	Create(rs app.RequestScope, company *models.Companies) error
	// Update updates the company with given ID in the storage.
	Update(rs app.RequestScope, id int, company *models.Companies) error
	// Delete removes the company with given ID from the storage.
	Delete(rs app.RequestScope, id int) error
}

// CompanyService provides services related with companys.
type CompanyService struct {
	dao companyDAO
}

// NewCompanyService creates a new CompanyService with the given company DAO.
func NewCompanyService(dao companyDAO) *CompanyService {
	return &CompanyService{dao}
}

// Get returns the company with the specified the company ID.
func (s *CompanyService) Get(rs app.RequestScope, id int) (*models.Companies, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new company.
func (s *CompanyService) Create(rs app.RequestScope, model *models.Companies) (*models.Companies, error) {
	if err := model.ValidateCompanies(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.CompanyID)
}

// Update updates the company with the specified ID.
func (s *CompanyService) Update(rs app.RequestScope, id int, model *models.Companies) (*models.Companies, error) {
	if err := model.ValidateCompanies(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the company with the specified ID.
func (s *CompanyService) Delete(rs app.RequestScope, id int) (*models.Companies, error) {
	company, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return company, err
}

// Count returns the number of companys.
func (s *CompanyService) Count(rs app.RequestScope) (int, error) {
	return s.dao.Count(rs)
}

// Query returns the companys with the specified offset and limit.
func (s *CompanyService) Query(rs app.RequestScope, offset, limit int) ([]models.Companies, error) {
	return s.dao.Query(rs, offset, limit)
}
