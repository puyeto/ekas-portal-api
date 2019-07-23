package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// DiscountDAO persists discount data in database
type DiscountDAO struct{}

// NewDiscountDAO creates a new DiscountDAO
func NewDiscountDAO() *DiscountDAO {
	return &DiscountDAO{}
}

// GetDiscount ...
func (dao *DiscountDAO) GetDiscount(rs app.RequestScope, id int32) (*models.Discounts, error) {
	var dis models.Discounts
	err := rs.Tx().Select().Model(id, &dis)
	return &dis, err
}

// GetDiscounts ...
func (dao *DiscountDAO) GetDiscounts(rs app.RequestScope, offset, limit int) ([]models.Discounts, error) {
	dis := []models.Discounts{}
	err := rs.Tx().Select().OrderBy("discount_id").Offset(int64(offset)).Limit(int64(limit)).All(&dis)
	return dis, err
}

// AddDiscounts ...
func (dao *DiscountDAO) AddDiscounts(rs app.RequestScope, dis *models.Discounts) error {
	dis.DiscountID = 0
	return rs.Tx().Model(dis).Insert()
}

// CountDiscounts returns the number of the discount records in the database.
func (dao *DiscountDAO) CountDiscounts(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("discounts").Row(&count)
	return count, err
}
