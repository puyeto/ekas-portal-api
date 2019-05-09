package apis

import (
	"fmt"
	"strings"
	// "time"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// vehicleService specifies the interface for the vehicle service needed by vehicleResource.
	vehicleService interface {
		GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error)
		GetTripDataByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.TripData, error)
		FetchAllTripsBetweenDates(rs app.RequestScope, deviceid string, offset, limit int, from string, to string) ([]models.TripData, error)
		Create(rs app.RequestScope, model *models.Vehicle) (int, error)
		CountTripRecordsBtwDates(rs app.RequestScope, deviceid string, from string, to string) (int, error)
		CountTripRecords(rs app.RequestScope, deviceid string) (int, error)
	}

	// vehicleResource defines the handlers for the CRUD APIs.
	vehicleResource struct {
		service vehicleService
	}
)

// ServeVehicleResource sets up the routing of vehicle endpoints and the corresponding handlers.
func ServeVehicleResource(rg *routing.RouteGroup, service vehicleService) {
	r := &vehicleResource{service}
	rg.Post("/addvehicle", r.create)
	rg.Get("/getconfigdetailsbystrid/<id>", r.getConfigurationByStringID)
	rg.Get("/gettripdata/<id>", r.getTripDataByDeviceID)
	rg.Post("/gettripdatabtwdates", r.getTripDataByDeviceIDBtwDates)
}

func (r *vehicleResource) getConfigurationByStringID(c *routing.Context) error {
	id := strings.ToLower(c.Param("id"))

	response, err := r.service.GetVehicleByStrID(app.GetRequestScope(c), id)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *vehicleResource) create(c *routing.Context) error {
	var model models.Vehicle
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.Create(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

// getTripDataByDeviceID ...
func (r *vehicleResource) getTripDataByDeviceID(c *routing.Context) error {
	deviceid := c.Param("id")

	rs := app.GetRequestScope(c)
	count, err := r.service.CountTripRecords(rs, deviceid)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	response, err := r.service.GetTripDataByDeviceID(rs, deviceid, paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

// getTripDataByDeviceID ...
func (r *vehicleResource) getTripDataByDeviceIDBtwDates(c *routing.Context) error {
	var model models.TripBetweenDates
	if err := c.Read(&model); err != nil {
		return err
	}

	// layout := "2006-01-02 11:45:26"	
	// formatedFrom, _ := time.Parse("Y-m-d H:i:s", model.From)
	// formatedTo, _ := time.Parse("Y-m-d H:i:s", model.To)

	fmt.Println(model)
	fmt.Println(model.From)

	rs := app.GetRequestScope(c)
	count, err := r.service.CountTripRecordsBtwDates(rs, model.DeviceID, model.From, model.To)
	if err != nil {
		return err
	}

	paginatedList := getPaginatedListFromRequest(c, count)
	if model.Offset == 0 {
		model.Offset = paginatedList.Offset()
	}
	if model.Limit == 0 {
		model.Limit = paginatedList.Limit()
	}
	response, err := r.service.FetchAllTripsBetweenDates(rs, model.DeviceID, model.Offset, model.Limit, model.From, model.To)
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}
