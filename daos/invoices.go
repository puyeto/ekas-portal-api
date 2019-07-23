package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// InvoiceDAO persists invoice data in database
type InvoiceDAO struct{}

// NewInvoiceDAO creates a new InvoiceDAO
func NewInvoiceDAO() *InvoiceDAO {
	return &InvoiceDAO{}
}

// GetInvoice ...
func (dao *InvoiceDAO) GetInvoice(rs app.RequestScope, id int32) (*models.Invoices, error) {
	var dis models.Invoices
	err := rs.Tx().Select().Model(id, &dis)
	return &dis, err
}

// GetInvoices ...
func (dao *InvoiceDAO) GetInvoices(rs app.RequestScope, offset, limit int) ([]models.Invoices, error) {
	dis := []models.Invoices{}
	err := rs.Tx().Select().OrderBy("invoice_id").Offset(int64(offset)).Limit(int64(limit)).All(&dis)
	return dis, err
}

// AddInvoices ...
func (dao *InvoiceDAO) AddInvoices(rs app.RequestScope, dis *models.Invoices) error {
	dis.InvoiceID = 0
	return rs.Tx().Model(dis).Insert()
}

// CountInvoices returns the number of the invoice records in the database.
func (dao *InvoiceDAO) CountInvoices(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("invoices").Row(&count)
	return count, err
}

// GetPricing ...
func (dao *InvoiceDAO) GetPricing(rs app.RequestScope, id int32) (*models.Pricings, error) {
	var model models.Pricings
	err := rs.Tx().Select().Model(id, &model)
	return &model, err
}

// GetDiscount ...
func (dao *InvoiceDAO) GetDiscount(rs app.RequestScope, id int32) (*models.Discounts, error) {
	var model models.Discounts
	err := rs.Tx().Select().Model(id, &model)
	return &model, err
}

// GetDiscount ...
func (dao *InvoiceDAO) GetVehicle(rs app.RequestScope, id int32) (*models.VehicleDetails, error) {
	var model models.VehicleDetails
	err := rs.Tx().Select().Model(id, &model)
	return &model, err
}
