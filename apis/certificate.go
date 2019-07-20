package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// certificateService specifies the interface for the certificate service needed by certificateResource.
	certificateService interface {
		Get(rs app.RequestScope, id int) (*models.Certificates, error)
		Query(rs app.RequestScope, offset, limit int) ([]models.Certificates, error)
		Count(rs app.RequestScope) (int, error)
		Create(rs app.RequestScope, model *models.Certificates) (*models.Certificates, error)
		Update(rs app.RequestScope, id int, model *models.Certificates) (*models.Certificates, error)
		Delete(rs app.RequestScope, id int) (*models.Certificates, error)
	}

	// certificateResource defines the handlers for the CRUD APIs.
	certificateResource struct {
		service certificateService
	}
)

// ServeCertificateResource sets up the routing of certificate endpoints and the corresponding handlers.
func ServeCertificateResource(rg *routing.RouteGroup, service certificateService) {
	r := &certificateResource{service}
	rg.Get("/certificates/get/<id>", r.get)
	rg.Get("/certificates/list", r.query)
	rg.Post("/certificates/create", r.create)
	rg.Put("/certificates/update/<id>", r.update)
	rg.Delete("/certificates/delete/<id>", r.delete)
}

func (r *certificateResource) get(c *routing.Context) error {
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

func (r *certificateResource) query(c *routing.Context) error {
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

func (r *certificateResource) create(c *routing.Context) error {
	var model models.Certificates
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *certificateResource) update(c *routing.Context) error {
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

func (r *certificateResource) delete(c *routing.Context) error {
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
