package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type (
	// saccoService specifies the interface for the sacco service needed by saccoResource.
	saccoService interface {
		Get(rs app.RequestScope, id int) (*models.Saccos, error)
		GetSaccoUser(rs app.RequestScope, userid int) (*models.Saccos, error)
		Query(rs app.RequestScope, offset, limit int) ([]models.Saccos, error)
		Count(rs app.RequestScope) (int, error)
		Create(rs app.RequestScope, model *models.Saccos) (*models.Saccos, error)
		Update(rs app.RequestScope, id int, model *models.Saccos) (*models.Saccos, error)
		Delete(rs app.RequestScope, id int) (*models.Saccos, error)
	}

	// saccoResource defines the handlers for the CRUD APIs.
	saccoResource struct {
		service saccoService
	}
)

// ServeSaccoResource sets up the routing of sacco endpoints and the corresponding handlers.
func ServeSaccoResource(rg *routing.RouteGroup, service saccoService) {
	r := &saccoResource{service}
	rg.Get("/saccos/get/<id>", r.get)
	rg.Get("/saccos/user/<id>", r.getSaccoUser)
	rg.Get("/saccos/list", r.query)
	rg.Post("/saccos/create", r.create)
	rg.Put("/sacco/update/<id>", r.update)
	rg.Delete("/saccos/del/<id>", r.delete)
}

func (r *saccoResource) get(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Get(app.GetRequestScope(c), id)
	if err != nil {
		return err
	}

	return c.Write(response)
}

// get sacco Details of a user
func (r *saccoResource) getSaccoUser(c *routing.Context) error {
	userid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.GetSaccoUser(app.GetRequestScope(c), userid)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *saccoResource) query(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	count, err := r.service.Count(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.Query(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *saccoResource) create(c *routing.Context) error {
	var model models.Saccos
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *saccoResource) update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)

	model, err := r.service.Get(rs, id)
	if err != nil {
		return err
	}

	if err := c.Read(model); err != nil {
		return err
	}

	response, err := r.service.Update(rs, id, model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *saccoResource) delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Delete(app.GetRequestScope(c), id)
	if err != nil {
		return err
	}

	return c.Write(response)
}
