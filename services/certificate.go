package services

import (
	"time"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// certificateDAO specifies the interface of the certificate DAO needed by CertificateService.
type certificateDAO interface {
	// Get returns the certificate with the specified certificate ID.
	Get(rs app.RequestScope, id int) (*models.Certificates, error)
	// Count returns the number of certificates.
	Count(rs app.RequestScope) (int, error)
	// Query returns the list of certificates with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]models.Certificates, error)
	// Create saves a new certificate in the storage.
	Create(rs app.RequestScope, certificate *models.Certificates) error
	// Update updates the certificate with given ID in the storage.
	Update(rs app.RequestScope, id int, certificate *models.Certificates) error
	// Delete removes the certificate with given ID from the storage.
	Delete(rs app.RequestScope, id int) error
}

// CertificateService provides services related with certificates.
type CertificateService struct {
	dao certificateDAO
}

// NewCertificateService creates a new CertificateService with the given certificate DAO.
func NewCertificateService(dao certificateDAO) *CertificateService {
	return &CertificateService{dao}
}

// Get returns the certificate with the specified the certificate ID.
func (s *CertificateService) Get(rs app.RequestScope, id int) (*models.Certificates, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new certificate.
func (s *CertificateService) Create(rs app.RequestScope, model *models.Certificates) (*models.Certificates, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}

	if model.IssuedOn.IsZero() {
		model.IssuedOn = time.Now()
	}

	if model.CertSerial == "" {
		model.CertSerial = model.CertNo
	}

	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.ID)
}

// Update updates the certificate with the specified ID.
func (s *CertificateService) Update(rs app.RequestScope, id int, model *models.Certificates) (*models.Certificates, error) {
	if err := model.Validate(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the certificate with the specified ID.
func (s *CertificateService) Delete(rs app.RequestScope, id int) (*models.Certificates, error) {
	certificate, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return certificate, err
}

// Count returns the number of certificates.
func (s *CertificateService) Count(rs app.RequestScope) (int, error) {
	return s.dao.Count(rs)
}

// Query returns the certificates with the specified offset and limit.
func (s *CertificateService) Query(rs app.RequestScope, offset, limit int) ([]models.Certificates, error) {
	return s.dao.Query(rs, offset, limit)
}
