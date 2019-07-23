package services

import (
	"errors"
	"time"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// invoiceDAO specifies the interface of the invoice DAO needed by InvoiceService.
type invoiceDAO interface {
	GetInvoice(rs app.RequestScope, id int32) (*models.Invoices, error)
	GetInvoices(rs app.RequestScope, offset, limit int) ([]models.Invoices, error)
	AddInvoices(rs app.RequestScope, model *models.Invoices) error
	CountInvoices(rs app.RequestScope) (int, error)
	GetPricing(rs app.RequestScope, id int32) (*models.Pricings, error)
	GetDiscount(rs app.RequestScope, id int32) (*models.Discounts, error)
	GetVehicle(rs app.RequestScope, id int32) (*models.VehicleDetails, error)
}

// InvoiceService provides services related with invoices.
type InvoiceService struct {
	dao invoiceDAO
}

// NewInvoiceService creates a new InvoiceService with the given invoice DAO.
func NewInvoiceService(dao invoiceDAO) *InvoiceService {
	return &InvoiceService{dao}
}

// GetInvoice returns the invoice with the specified the artist ID.
func (s *InvoiceService) GetInvoice(rs app.RequestScope, id int32) (*models.Invoices, error) {
	return s.dao.GetInvoice(rs, id)
}

// GetInvoices get alist of invoices.
func (s *InvoiceService) GetInvoices(rs app.RequestScope, offset, limit int) ([]models.Invoices, error) {
	return s.dao.GetInvoices(rs, offset, limit)
}

// AddInvoices ...
func (s *InvoiceService) AddInvoices(rs app.RequestScope, model *models.Invoices) (*models.Invoices, error) {

	m, err := s.newInvoice(rs, model)
	if err != nil {
		return nil, err
	}

	if err := m.ValidateInvoice(); err != nil {
		return nil, err
	}
	if err := s.dao.AddInvoices(rs, m); err != nil {
		return nil, err
	}
	return s.dao.GetInvoice(rs, m.InvoiceID)
}

func (s *InvoiceService) newInvoice(rs app.RequestScope, model *models.Invoices) (*models.Invoices, error) {
	// get vehicle last invoice duedate
	vehicle, err := s.dao.GetVehicle(rs, model.VehicleID)
	if err != nil {
		return nil, err
	}

	if vehicle.InvoicingStatus == 0 {
		return nil, errors.New("Automatic invoicing of this vehicle is disabled")
	}

	// get invoice amount
	price, err := s.dao.GetPricing(rs, model.PricingID)
	if err != nil {
		return nil, err
	}

	// get discount amount
	discount, err := s.dao.GetDiscount(rs, model.DiscountID)
	if err != nil {
		return nil, err
	}

	model.InvoiceAmount = price.PricingAmount
	model.InvoiceTotal = price.PricingAmount - discount.DiscountValue
	if discount.DiscountValueAs == "percent" {
		model.InvoiceTotal = price.PricingAmount - (discount.DiscountValue / 100 * price.PricingAmount)
	}

	date := time.Now()
	invDate, _ := time.Parse(time.RFC3339, date.Format(time.RFC3339))
	model.InvoicingDate = invDate

	dueDate, _ := time.Parse(time.RFC3339, vehicle.InvoiceDueDate.Format(time.RFC3339))
	model.InvoiceStartDate = dueDate.AddDate(0, 0, 1)

	diff := invDate.Sub(dueDate)
	if int(diff.Hours()/24) > 0 {
		model.InvoiceStartDate = invDate
	}

	switch price.PricingDuration {
	case "Daily":
		model.InvoiceEndDate = model.InvoiceStartDate.AddDate(0, 0, 1)
	case "Weekly":
		model.InvoiceEndDate = model.InvoiceStartDate.AddDate(0, 0, 7)
	case "Monthly":
		model.InvoiceEndDate = model.InvoiceStartDate.AddDate(0, 1, 0)
	case "Bi-Monthly":
		model.InvoiceEndDate = model.InvoiceStartDate.AddDate(0, 2, 0)
	case "Quarterly":
		model.InvoiceEndDate = model.InvoiceStartDate.AddDate(0, 3, 0)
	case "Semi-Annually":
		model.InvoiceEndDate = model.InvoiceStartDate.AddDate(0, 6, 0)
	case "Annually":
		model.InvoiceEndDate = model.InvoiceStartDate.AddDate(1, 0, 0)
	default:
		model.InvoiceEndDate = model.InvoiceStartDate.AddDate(0, 1, 0)
	}

	return model, err
}

// CountInvoices returns the number of artists.
func (s *InvoiceService) CountInvoices(rs app.RequestScope) (int, error) {
	return s.dao.CountInvoices(rs)
}
