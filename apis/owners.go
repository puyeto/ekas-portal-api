package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// ownerService specifies the interface for the owner service needed by ownerResource.
	ownerService interface {
		Get(rs app.RequestScope, id uint32) (*models.VehicleOwner, error)
		Query(rs app.RequestScope, offset, limit, uid int) ([]models.VehicleOwner, error)
		Count(rs app.RequestScope, uid int) (int, error)
		Create(rs app.RequestScope, model *models.VehicleOwner) (*models.VehicleOwner, error)
		Update(rs app.RequestScope, id uint32, model *models.VehicleOwner) (*models.VehicleOwner, error)
		Delete(rs app.RequestScope, id uint32) (*models.VehicleOwner, error)
	}

	// ownerResource defines the handlers for the CRUD APIs.
	ownerResource struct {
		service ownerService
	}
)

// ServeOwner sets up the routing of owner endpoints and the corresponding handlers.
func ServeOwnerResource(rg *routing.RouteGroup, service ownerService) {
	r := &ownerResource{service}
	rg.Get("/owner/get/<id>", r.get)
	rg.Get("/owners/list", r.query)
	rg.Post("/owners/create", r.create)
	rg.Put("/owners/update/<id>", r.update)
	rg.Delete("/owners/del/<id>", r.delete)
}

func (r *ownerResource) get(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Get(app.GetRequestScope(c), uint32(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *ownerResource) query(c *routing.Context) error {
	// get company id from query string
	uid, _ := strconv.Atoi(c.Query("uid", "0"))

	rs := app.GetRequestScope(c)
	count, err := r.service.Count(rs, uid)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.Query(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit(), uid)
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *ownerResource) create(c *routing.Context) error {
	var model models.VehicleOwner
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *ownerResource) update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)

	model, err := r.service.Get(rs, uint32(id))
	if err != nil {
		return err
	}

	if err := c.Read(model); err != nil {
		return err
	}

	response, err := r.service.Update(rs, uint32(id), model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *ownerResource) delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Delete(app.GetRequestScope(c), uint32(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}
