package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// invoiceService specifies the interface for the invoice service needed by invoiceResource.
	invoiceService interface {
		GetInvoice(rs app.RequestScope, id int32) (*models.Invoices, error)
		GetInvoices(rs app.RequestScope, offset, limit int) ([]models.Invoices, error)
		AddInvoices(rs app.RequestScope, model *models.Invoices) (*models.Invoices, error)
		CountInvoices(rs app.RequestScope) (int, error)
	}

	// invoiceResource defines the handlers for the CRUD APIs.
	invoiceResource struct {
		service invoiceService
	}
)

// ServeInvoiceResource sets up the routing of invoice endpoints and the corresponding handlers.
func ServeInvoiceResource(rg *routing.RouteGroup, service invoiceService) {
	r := &invoiceResource{service}
	rg.Get("/invoice/<id>", r.getinvoice)
	rg.Get("/invoices", r.getinvoices)
	rg.Post("/invoice", r.addinvoices)
}

func (r *invoiceResource) getinvoices(c *routing.Context) error {

	rs := app.GetRequestScope(c)
	count, err := r.service.CountInvoices(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.GetInvoices(rs, paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *invoiceResource) getinvoice(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.GetInvoice(app.GetRequestScope(c), int32(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *invoiceResource) addinvoices(c *routing.Context) error {
	var model models.Invoices
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.AddInvoices(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}
