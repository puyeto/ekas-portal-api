package apis

import (
	"fmt"
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type (
	// deviceService specifies the interface for the device service needed by deviceResource.
	deviceService interface {
		Get(rs app.RequestScope, id int32) (*models.Devices, error)
		Query(rs app.RequestScope, offset, limit, cid int) ([]models.Devices, error)
		QueryPositions(rs app.RequestScope, offset, limit int, uid uint32, start, stop int64) ([]models.Devices, error)
		CountQueryPositions(rs app.RequestScope, uid uint32) (int, error)
		Count(rs app.RequestScope, cid int) (int, error)
		Create(rs app.RequestScope, model *models.Devices) (*models.Devices, error)
		Update(rs app.RequestScope, id int32, model *models.Devices) (*models.Devices, error)
		Delete(rs app.RequestScope, id int32) (*models.Devices, error)
		CountConfiguredDevices(rs app.RequestScope, vehicleid int, deviceid int64) (int, error)
		ConfiguredDevices(rs app.RequestScope, offset, limit, vehicleid int, deviceid int64) ([]models.DeviceConfiguration, error)
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
	rg.Get("/devices/positions", r.queryPositions)
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
	// get company id from query string
	cid, _ := strconv.Atoi(c.Query("cid", "0"))

	rs := app.GetRequestScope(c)
	count, err := r.service.Count(rs, cid)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.Query(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit(), cid)
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *deviceResource) queryPositions(c *routing.Context) error {
	uid, err := strconv.Atoi(c.Query("uid", "0"))
	if err != nil {
		return err
	}

	start, err := strconv.Atoi(c.Query("start", "0"))
	if err != nil {
		return err
	}

	stop, err := strconv.Atoi(c.Query("stop", "99"))
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)
	count, err := r.service.CountQueryPositions(rs, uint32(uid))
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.QueryPositions(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit(), uint32(uid), int64(start), int64(stop))
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *deviceResource) count(c *routing.Context) error {
	// get company id from query string
	cid, _ := strconv.Atoi(c.Query("cid", "0"))

	response, err := r.service.Count(app.GetRequestScope(c), cid)
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
	vehicleid, err := strconv.Atoi(c.Query("vehicle_id", "0"))
	if err != nil {
		return err
	}

	deviceid, err := strconv.ParseInt(c.Query("device_id", "0"), 10, 64)
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)
	count, err := r.service.CountConfiguredDevices(rs, vehicleid, deviceid)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.ConfiguredDevices(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit(), vehicleid, deviceid)
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}
