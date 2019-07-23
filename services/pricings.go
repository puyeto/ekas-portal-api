package services

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// pricingDAO specifies the interface of the pricing DAO needed by PricingService.
type pricingDAO interface {
	GetPricing(rs app.RequestScope, id int32) (*models.Pricings, error)
	GetPricings(rs app.RequestScope, offset, limit int) ([]models.Pricings, error)
	AddPricings(rs app.RequestScope, model *models.Pricings) error
	CountPricings(rs app.RequestScope) (int, error)
}

// PricingService provides services related with pricings.
type PricingService struct {
	dao pricingDAO
}

// NewPricingService creates a new PricingService with the given pricing DAO.
func NewPricingService(dao pricingDAO) *PricingService {
	return &PricingService{dao}
}

// GetPricing returns the pricing with the specified the artist ID.
func (s *PricingService) GetPricing(rs app.RequestScope, id int32) (*models.Pricings, error) {
	return s.dao.GetPricing(rs, id)
}

// GetPricings get alist of pricings.
func (s *PricingService) GetPricings(rs app.RequestScope, offset, limit int) ([]models.Pricings, error) {
	return s.dao.GetPricings(rs, offset, limit)
}

// AddPricings ...
func (s *PricingService) AddPricings(rs app.RequestScope, model *models.Pricings) (*models.Pricings, error) {
	if err := model.ValidatePricing(); err != nil {
		return nil, err
	}
	if err := s.dao.AddPricings(rs, model); err != nil {
		return nil, err
	}
	return s.dao.GetPricing(rs, model.PricingID)
}

// Count returns the number of artists.
func (s *PricingService) CountPricings(rs app.RequestScope) (int, error) {
	return s.dao.CountPricings(rs)
}
