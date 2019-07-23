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
		ListInvoices(rs app.RequestScope, offset, limit, uid, vid int) ([]models.Invoices, error)
		AddInvoices(rs app.RequestScope, model *models.Invoices) (*models.Invoices, error)
		CountInvoices(rs app.RequestScope, uid, vid int) (int, error)
	}

	// invoiceResource defines the handlers for the CRUD APIs.
	invoiceResource struct {
		service invoiceService
	}
)

// ServeInvoiceResource sets up the routing of invoice endpoints and the corresponding handlers.
func ServeInvoiceResource(rg *routing.RouteGroup, service invoiceService) {
	r := &invoiceResource{service}
	rg.Get("/invoice/get/<id>", r.getinvoice)
	rg.Get("/invoices/list", r.listinvoices)
	rg.Post("/invoice/create", r.addinvoices)
}

func (r *invoiceResource) listinvoices(c *routing.Context) error {
	uid, err := strconv.Atoi(c.Query("uid", "0"))
	if err != nil {
		return err
	}
	vid, err := strconv.Atoi(c.Query("vid", "0"))
	if err != nil {
		return err
	}
	rs := app.GetRequestScope(c)
	count, err := r.service.CountInvoices(rs, uid, vid)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.ListInvoices(rs, paginatedList.Offset(), paginatedList.Limit(), uid, vid)
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
