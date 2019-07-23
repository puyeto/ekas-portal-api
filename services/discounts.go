package services

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// discountDAO specifies the interface of the discount DAO needed by DiscountService.
type discountDAO interface {
	GetDiscount(rs app.RequestScope, id int32) (*models.Discounts, error)
	GetDiscounts(rs app.RequestScope, offset, limit int) ([]models.Discounts, error)
	AddDiscounts(rs app.RequestScope, model *models.Discounts) error
	CountDiscounts(rs app.RequestScope) (int, error)
}

// DiscountService provides services related with discounts.
type DiscountService struct {
	dao discountDAO
}

// NewDiscountService creates a new DiscountService with the given discount DAO.
func NewDiscountService(dao discountDAO) *DiscountService {
	return &DiscountService{dao}
}

// GetDiscount returns the discount with the specified the artist ID.
func (s *DiscountService) GetDiscount(rs app.RequestScope, id int32) (*models.Discounts, error) {
	return s.dao.GetDiscount(rs, id)
}

// GetDiscounts get alist of discounts.
func (s *DiscountService) GetDiscounts(rs app.RequestScope, offset, limit int) ([]models.Discounts, error) {
	return s.dao.GetDiscounts(rs, offset, limit)
}

// AddDiscounts ...
func (s *DiscountService) AddDiscounts(rs app.RequestScope, model *models.Discounts) (*models.Discounts, error) {
	if err := model.ValidateDiscount(); err != nil {
		return nil, err
	}
	if err := s.dao.AddDiscounts(rs, model); err != nil {
		return nil, err
	}
	return s.dao.GetDiscount(rs, model.DiscountID)
}

// Count returns the number of artists.
func (s *DiscountService) CountDiscounts(rs app.RequestScope) (int, error) {
	return s.dao.CountDiscounts(rs)
}
