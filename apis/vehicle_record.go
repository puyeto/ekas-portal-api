package apis

import (
	"fmt"
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type (
	// vehicleRecordService specifies the interface for the vehicleRecord service needed by vehicleRecordResource.
	vehicleRecordService interface {
		Get(rs app.RequestScope, id uint32) (*models.VehicleDetails, error)
		CreateVehicle(rs app.RequestScope, model *models.VehicleDetails) (uint32, error)
		Query(rs app.RequestScope, offset, limit int, uid int, typ string, userdetails models.AuthUsers) ([]models.VehicleDetails, error)
		QueryFilter(rs app.RequestScope, offset, limit int, m *models.FilterVehicles) ([]models.VehicleDetails, error)
		Count(rs app.RequestScope, uid int, typ string, userdetails models.AuthUsers) (int, error)
		CountFilter(rs app.RequestScope, m *models.FilterVehicles) (int, error)
		Update(rs app.RequestScope, model *models.VehicleDetails) error
		Delete(rs app.RequestScope, id uint32) error
		RenewVehicle(rs app.RequestScope, model *models.VehicleRenewals) (uint32, error)
		ListVehicleRenewals(rs app.RequestScope, offset, limit int) ([]models.VehicleRenewals, error)
		CountRenewals(rs app.RequestScope) (int, error)
		CreateReminder(rs app.RequestScope, model *models.Reminders) (uint32, error)
		CountReminders(rs app.RequestScope, uid int) (int, error)
		GetReminder(rs app.RequestScope, offset, limit int, uid int) ([]models.Reminders, error)
		GetUser(rs app.RequestScope, id uint32) (models.AuthUsers, error)
	}

	// vehicleRecordResource defines the handlers for the CRUD APIs.
	vehicleRecordResource struct {
		service vehicleRecordService
	}
)

// ServeVehicleRecordResource sets up the routing of vehicleRecord endpoints and the corresponding handlers.
func ServeVehicleRecordResource(rg *routing.RouteGroup, service vehicleRecordService) {
	r := &vehicleRecordResource{service}
	rg.Get("/vehicle/get/<id>", r.get)
	rg.Post("/vehicle/create", r.createVehicle)
	rg.Get("/vehicles/list", r.query)
	rg.Post("/vehicles/filter", r.filterVehicles)
	rg.Get("/vehicles/count", r.count)
	rg.Put("/vehicle/update", r.update)
	rg.Delete("/vehicle/delete/<id>", r.delete)
	rg.Post("/vehicle/renew", r.renewVehicle)
	rg.Get("/vehicle/renew", r.listVehicleRenewals)
	rg.Post("/vehicle/reminders", r.createReminder)
	rg.Get("/vehicle/reminders", r.getReminder)
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

func (r *vehicleRecordResource) createVehicle(c *routing.Context) error {
	var model models.VehicleDetails
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.CreateVehicle(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *vehicleRecordResource) query(c *routing.Context) error {
	uid, err := strconv.Atoi(c.Query("uid", "0"))
	if err != nil {
		return err
	}

	// get user details
	userdetails, _ := r.service.GetUser(app.GetRequestScope(c), uint32(uid))
	fmt.Println(userdetails)

	typ := c.Query("type", "")

	rs := app.GetRequestScope(c)
	count, err := r.service.Count(rs, uid, typ, userdetails)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.Query(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit(), uid, typ, userdetails)
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *vehicleRecordResource) filterVehicles(c *routing.Context) error {
	var model models.FilterVehicles
	if err := c.Read(&model); err != nil {
		return err
	}

	rs := app.GetRequestScope(c)
	count, err := r.service.CountFilter(rs, &model)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.QueryFilter(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit(), &model)
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *vehicleRecordResource) update(c *routing.Context) error {
	var model models.VehicleDetails
	if err := c.Read(&model); err != nil {
		return err
	}

	err := r.service.Update(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(map[string]uint32{
		"updated_id": model.VehicleID,
	})
}

func (r *vehicleRecordResource) delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	err = r.service.Delete(app.GetRequestScope(c), uint32(id))
	if err != nil {
		return err
	}

	return c.Write(map[string]int{
		"deleted_id": id,
	})
}

func (r *vehicleRecordResource) count(c *routing.Context) error {
	typ := c.Query("type", "")
	uid, err := strconv.Atoi(c.Query("uid", "0"))
	if err != nil {
		return err
	}

	// get user details
	userdetails, _ := r.service.GetUser(app.GetRequestScope(c), uint32(uid))
	fmt.Println(userdetails)

	response, err := r.service.Count(app.GetRequestScope(c), uid, typ, userdetails)
	if err != nil {
		return err
	}

	return c.Write(map[string]int{
		"count": response,
	})
}

func (r *vehicleRecordResource) renewVehicle(c *routing.Context) error {
	var model models.VehicleRenewals
	if err := c.Read(&model); err != nil {
		return err
	}

	response, err := r.service.RenewVehicle(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(map[string]uint32{
		"response": response,
	})
}

func (r *vehicleRecordResource) listVehicleRenewals(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	count, err := r.service.CountRenewals(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.ListVehicleRenewals(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *vehicleRecordResource) createReminder(c *routing.Context) error {
	var model models.Reminders
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.CreateReminder(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(map[string]uint32{
		"response": response,
	})
}

func (r *vehicleRecordResource) getReminder(c *routing.Context) error {
	uid, err := strconv.Atoi(c.Query("uid", "0"))
	if err != nil {
		return err
	}

	rs := app.GetRequestScope(c)
	count, err := r.service.CountReminders(rs, uid)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.GetReminder(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit(), uid)
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}
