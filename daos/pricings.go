package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// PricingDAO persists pricing data in database
type PricingDAO struct{}

// NewPricingDAO creates a new PricingDAO
func NewPricingDAO() *PricingDAO {
	return &PricingDAO{}
}

// GetPricing ...
func (dao *PricingDAO) GetPricing(rs app.RequestScope, id int32) (*models.Pricings, error) {
	var dis models.Pricings
	err := rs.Tx().Select().Model(id, &dis)
	return &dis, err
}

// GetPricings ...
func (dao *PricingDAO) GetPricings(rs app.RequestScope, offset, limit int) ([]models.Pricings, error) {
	dis := []models.Pricings{}
	err := rs.Tx().Select().OrderBy("pricing_id").Offset(int64(offset)).Limit(int64(limit)).All(&dis)
	return dis, err
}

// AddPricings ...
func (dao *PricingDAO) AddPricings(rs app.RequestScope, dis *models.Pricings) error {
	dis.PricingID = 0
	return rs.Tx().Model(dis).Insert()
}

// CountPricings returns the number of the pricing records in the database.
func (dao *PricingDAO) CountPricings(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("pricings").Row(&count)
	return count, err
}
