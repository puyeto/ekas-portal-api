package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type (
	// simcardService specifies the interface for the simcard service needed by simcardResource.
	simcardService interface {
		Get(rs app.RequestScope, id int) (*models.Simcards, error)
		Query(rs app.RequestScope, offset, limit int) ([]models.Simcards, error)
		Count(rs app.RequestScope) (int, error)
		Create(rs app.RequestScope, model *models.Simcards) (*models.Simcards, error)
		Update(rs app.RequestScope, id int, model *models.Simcards) (*models.Simcards, error)
		Delete(rs app.RequestScope, id int) (*models.Simcards, error)
		GetStats(rs app.RequestScope) *models.SimcardStats
	}

	// simcardResource defines the handlers for the CRUD APIs.
	simcardResource struct {
		service simcardService
	}
)

// ServeSimcardResource sets up the routing of simcard endpoints and the corresponding handlers.
func ServeSimcardResource(rg *routing.RouteGroup, service simcardService) {
	r := &simcardResource{service}
	rg.Get("/simcards/get/<id>", r.get)
	rg.Get("/simcards/list", r.query)
	rg.Post("/simcards/create", r.create)
	rg.Put("/simcard/update/<id>", r.update)
	rg.Delete("/simcards/del/<id>", r.delete)
	rg.Get("/simcards/stats", r.stats)
}

func (r *simcardResource) get(c *routing.Context) error {
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

func (r *simcardResource) query(c *routing.Context) error {
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

func (r *simcardResource) create(c *routing.Context) error {
	var model models.Simcards
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *simcardResource) update(c *routing.Context) error {
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

func (r *simcardResource) delete(c *routing.Context) error {
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

func (r *simcardResource) stats(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	return c.Write(r.service.GetStats(rs))
}
