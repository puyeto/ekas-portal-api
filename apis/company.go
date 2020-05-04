package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type (
	// companyService specifies the interface for the company service needed by companyResource.
	companyService interface {
		Get(rs app.RequestScope, id int) (*models.Companies, error)
		GetCompanyUser(rs app.RequestScope, userid int) (*models.Companies, error)
		Query(rs app.RequestScope, offset, limit int) ([]models.Companies, error)
		Count(rs app.RequestScope) (int, error)
		Create(rs app.RequestScope, model *models.Companies) (*models.Companies, error)
		Update(rs app.RequestScope, id int, model *models.Companies) (*models.Companies, error)
		Delete(rs app.RequestScope, id int) (*models.Companies, error)
	}

	// companyResource defines the handlers for the CRUD APIs.
	companyResource struct {
		service companyService
	}
)

// ServeCompanyResource sets up the routing of company endpoints and the corresponding handlers.
func ServeCompanyResource(rg *routing.RouteGroup, service companyService) {
	r := &companyResource{service}
	rg.Get("/companies/get/<id>", r.get)
	rg.Get("/companies/user/<id>", r.getCompanyUser)
	rg.Get("/companies/list", r.query)
	rg.Post("/companies/create", r.create)
	rg.Put("/company/update/<id>", r.update)
	rg.Delete("/companies/del/<id>", r.delete)
}

func (r *companyResource) get(c *routing.Context) error {
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

// get company Details of a user
func (r *companyResource) getCompanyUser(c *routing.Context) error {
	userid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.GetCompanyUser(app.GetRequestScope(c), userid)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *companyResource) query(c *routing.Context) error {
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

func (r *companyResource) create(c *routing.Context) error {
	var model models.Companies
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *companyResource) update(c *routing.Context) error {
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

func (r *companyResource) delete(c *routing.Context) error {
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
