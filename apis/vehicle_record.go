package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// vehicleRecordService specifies the interface for the vehicleRecord service needed by vehicleRecordResource.
	vehicleRecordService interface {
		Get(rs app.RequestScope, id uint32) (*models.VehicleDetails, error)
		Query(rs app.RequestScope, offset, limit int) ([]models.VehicleDetails, error)
		Count(rs app.RequestScope) (int, error)
		Update(rs app.RequestScope, id uint32, model *models.VehicleDetails) (*models.VehicleDetails, error)
		Delete(rs app.RequestScope, id uint32) (*models.VehicleDetails, error)
	}

	// vehicleRecordResource defines the handlers for the CRUD APIs.
	vehicleRecordResource struct {
		service vehicleRecordService
	}
)

// ServeVehicleRecord sets up the routing of vehicleRecord endpoints and the corresponding handlers.
func ServeVehicleRecordResource(rg *routing.RouteGroup, service vehicleRecordService) {
	r := &vehicleRecordResource{service}
	rg.Get("/vehicleRecords/<id>", r.get)
	rg.Get("/vehicles/list", r.query)
	rg.Put("/vehicleRecords/<id>", r.update)
	rg.Delete("/vehicleRecords/<id>", r.delete)
}

func (r *vehicleRecordResource) get(c *routing.Context) error {
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

func (r *vehicleRecordResource) query(c *routing.Context) error {
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

func (r *vehicleRecordResource) update(c *routing.Context) error {
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

func (r *vehicleRecordResource) delete(c *routing.Context) error {
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
