package apis

import (
	"fmt"
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// deviceService specifies the interface for the device service needed by deviceResource.
	deviceService interface {
		Get(rs app.RequestScope, id int32) (*models.Devices, error)
		Query(rs app.RequestScope, offset, limit int) ([]models.Devices, error)
		Count(rs app.RequestScope) (int, error)
		Create(rs app.RequestScope, model *models.Devices) (*models.Devices, error)
		Update(rs app.RequestScope, id int32, model *models.Devices) (*models.Devices, error)
		Delete(rs app.RequestScope, id int32) (*models.Devices, error)
		CountConfiguredDevices(rs app.RequestScope) (int, error)
		ConfiguredDevices(rs app.RequestScope, offset, limit int) ([]models.DeviceConfiguration, error)
	}

	// deviceResource defines the handlers for the CRUD APIs.
	deviceResource struct {
		service deviceService
	}
)

// ServeDeviceResource sets up the routing of device endpoints and the corresponding handlers.
func ServeDeviceResource(rg *routing.RouteGroup, service deviceService) {
	r := &deviceResource{service}
	rg.Get("/device/<id>", r.get)
	rg.Get("/devices/list", r.query)
	rg.Get("/devices/count", r.count)
	rg.Get("/devices/configured-devices", r.configuredDevices)
	rg.Post("/devices/create", r.create)
	rg.Put("/devices/<id>", r.update)
	rg.Delete("/device/<id>", r.delete)
}

func (r *deviceResource) get(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Get(app.GetRequestScope(c), int32(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *deviceResource) query(c *routing.Context) error {
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

func (r *deviceResource) count(c *routing.Context) error {
	response, err := r.service.Count(app.GetRequestScope(c))
	if err != nil {
		return err
	}

	return c.Write(map[string]int{
		"count": response,
	})
}

func (r *deviceResource) create(c *routing.Context) error {
	var model models.Devices
	if err := c.Read(&model); err != nil {
		return err
	}
	fmt.Println(model)
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *deviceResource) update(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)

	model, err := r.service.Get(rs, int32(id))
	if err != nil {
		return err
	}

	if err := c.Read(model); err != nil {
		return err
	}

	response, err := r.service.Update(rs, int32(id), model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *deviceResource) delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.Delete(app.GetRequestScope(c), int32(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *deviceResource) configuredDevices(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	count, err := r.service.CountConfiguredDevices(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.ConfiguredDevices(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}
