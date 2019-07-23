package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// discountService specifies the interface for the discount service needed by discountResource.
	discountService interface {
		GetDiscount(rs app.RequestScope, id int32) (*models.Discounts, error)
		GetDiscounts(rs app.RequestScope, offset, limit int) ([]models.Discounts, error)
		AddDiscounts(rs app.RequestScope, model *models.Discounts) (*models.Discounts, error)
		CountDiscounts(rs app.RequestScope) (int, error)
	}

	// discountResource defines the handlers for the CRUD APIs.
	discountResource struct {
		service discountService
	}
)

// ServeDiscountResource sets up the routing of discount endpoints and the corresponding handlers.
func ServeDiscountResource(rg *routing.RouteGroup, service discountService) {
	r := &discountResource{service}
	rg.Get("/discount/<id>", r.getdiscount)
	rg.Get("/discounts", r.getdiscounts)
	rg.Post("/discount", r.adddiscounts)
}

func (r *discountResource) getdiscounts(c *routing.Context) error {

	rs := app.GetRequestScope(c)
	count, err := r.service.CountDiscounts(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.GetDiscounts(rs, paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *discountResource) getdiscount(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.GetDiscount(app.GetRequestScope(c), int32(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *discountResource) adddiscounts(c *routing.Context) error {
	var model models.Discounts
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.AddDiscounts(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}
