package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// pricingService specifies the interface for the pricing service needed by pricingResource.
	pricingService interface {
		GetPricing(rs app.RequestScope, id int32) (*models.Pricings, error)
		GetPricings(rs app.RequestScope, offset, limit int) ([]models.Pricings, error)
		AddPricings(rs app.RequestScope, model *models.Pricings) (*models.Pricings, error)
		CountPricings(rs app.RequestScope) (int, error)
	}

	// pricingResource defines the handlers for the CRUD APIs.
	pricingResource struct {
		service pricingService
	}
)

// ServePricingResource sets up the routing of pricing endpoints and the corresponding handlers.
func ServePricingResource(rg *routing.RouteGroup, service pricingService) {
	r := &pricingResource{service}
	rg.Get("/pricing/<id>", r.getpricing)
	rg.Get("/pricings", r.getpricings)
	rg.Post("/pricing", r.addpricings)
}

func (r *pricingResource) getpricings(c *routing.Context) error {

	rs := app.GetRequestScope(c)
	count, err := r.service.CountPricings(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.GetPricings(rs, paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *pricingResource) getpricing(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.GetPricing(app.GetRequestScope(c), int32(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *pricingResource) addpricings(c *routing.Context) error {
	var model models.Pricings
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.AddPricings(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}
