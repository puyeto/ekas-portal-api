package apis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// vehicleService specifies the interface for the vehicle service needed by vehicleResource.
	vehicleService interface {
		GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error)
		GetConfigurationDetails(rs app.RequestScope, vehicleid, deviceid int) (*models.VehicleConfigDetails, error)
		GetTripDataByDeviceID(deviceid string, offset, limit int) ([]models.DeviceData, error)
		FetchAllTripsBetweenDates(rs app.RequestScope, deviceid string, offset, limit int, from int64, to int64) ([]models.DeviceData, error)
		Create(rs app.RequestScope, model *models.Vehicle) (int, error)
		CountTripRecords(rs app.RequestScope, deviceid string) (int, error)
		CountRedisTripRecords(deviceid string) int
		CountRedisTripRecordsBtwDates(rs app.RequestScope, deviceid string, from int64, to int64) int
		CountOverspeed(rs app.RequestScope, deviceid string) (int, error)
		CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error)
		GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.TripData, error)
		GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.TripData, error)
		ListRecentViolations(rs app.RequestScope, offset, limit int, uid string) ([]models.CurrentViolations, error)
		GetCurrentViolations(rs app.RequestScope) ([]models.DeviceData, error)
		ListAllViolations(rs app.RequestScope, offset, limit int) ([]models.DeviceData, error)
		CountAllViolations(rs app.RequestScope) int
		SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int) ([]models.SearchDetails, error)
		CountSearches(rs app.RequestScope, searchterm string) (int, error)
		GetUnavailableDevices(rs app.RequestScope) ([]models.DeviceData, error)
		GetOfflineViolations(rs app.RequestScope, deviceid string) ([]models.DeviceData, error)
		CountTripDataByDeviceID(deviceid string) (int, error)
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
	rg.Post("/addvehicleconfiguration", r.create)
	rg.Get("/getconfigdetailsbystrid/<id>", r.getConfigurationByStringID)
	rg.Get("/device/data/<id>", r.getTripDataByDeviceID)
	rg.Get("/getoverspeed/<id>", r.getOverspeedsByDeviceID)
	rg.Get("/getfailsafe/<id>", r.getFailsafeByDeviceID)
	rg.Get("/getdisconnects/<id>", r.getDisconnectsByDeviceID)
	rg.Post("/gettripdatabtwdates", r.getTripDataByDeviceIDBtwDates)
	rg.Get("/listrecentviolations", r.listRecentViolations)
	rg.Get("/currentviolation", r.getCurrentViolations)
	rg.Get("/listviolations", r.listAllViolations)
	rg.Get("/getoffline/<id>", r.getOffline)
	rg.Get("/search/<term>", r.searchVehicle)
	rg.Get("/unavailable", r.getUnavailable)
	rg.Get("/configuration/details", r.getConfigurationDetails)

}

func (r *vehicleResource) getConfigurationDetails(c *routing.Context) error {
	// id := strings.ToLower(c.Param("id"))
	// get vehicle and deviceid from query string
	vid, _ := strconv.Atoi(c.Query("vid", "0"))
	did, _ := strconv.Atoi(c.Query("did", "0"))

	response, err := r.service.GetConfigurationDetails(app.GetRequestScope(c), vid, did)
	if err != nil {
		return err
	}

	return c.Write(response)
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

	// rs := app.GetRequestScope(c)
	count, _ := r.service.CountTripDataByDeviceID(deviceid)
	paginatedList := getPaginatedListFromRequest(c, count)

	response, err := r.service.GetTripDataByDeviceID(deviceid, paginatedList.Offset(), paginatedList.Limit())
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

	rs := app.GetRequestScope(c)
	count := r.service.CountRedisTripRecordsBtwDates(rs, model.DeviceID, model.From, model.To)
	paginatedList := getPaginatedListFromRequest(c, count)

	response, err := r.service.FetchAllTripsBetweenDates(rs, model.DeviceID, paginatedList.Offset(), paginatedList.Limit(), model.From, model.To)
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

// listAllViolations ...
func (r *vehicleResource) listAllViolations(c *routing.Context) error {

	rs := app.GetRequestScope(c)
	count := r.service.CountAllViolations(rs)
	paginatedList := getPaginatedListFromRequest(c, count)

	response, err := r.service.ListAllViolations(rs, paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

func (r *vehicleResource) getOffline(c *routing.Context) error {
	deviceid := c.Param("id")
	resp, err := r.service.GetOfflineViolations(app.GetRequestScope(c), deviceid)
	if err != nil {
		return err
	}
	return c.Write(resp)
}

// getOverspeedsByDeviceID ...
func (r *vehicleResource) getOverspeedsByDeviceID(c *routing.Context) error {
	deviceid := c.Param("id")

	rs := app.GetRequestScope(c)
	count, err := r.service.CountOverspeed(rs, deviceid)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	response, err := r.service.GetOverspeedByDeviceID(rs, deviceid, paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

// getFailsafeByDeviceID
func (r *vehicleResource) getFailsafeByDeviceID(c *routing.Context) error {
	deviceid := c.Param("id")

	rs := app.GetRequestScope(c)
	count, err := r.service.CountViolations(rs, deviceid, "failsafe")
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	response, err := r.service.GetViolationsByDeviceID(rs, deviceid, "failsafe", paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

// getDisconnectsByDeviceID
func (r *vehicleResource) getDisconnectsByDeviceID(c *routing.Context) error {
	deviceid := c.Param("id")

	rs := app.GetRequestScope(c)
	count, err := r.service.CountViolations(rs, deviceid, "disconnects")
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	response, err := r.service.GetViolationsByDeviceID(rs, deviceid, "disconnects", paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

// searchVehicle ...
func (r *vehicleResource) searchVehicle(c *routing.Context) error {
	searchterm := c.Param("term")
	rs := app.GetRequestScope(c)
	count, err := r.service.CountSearches(rs, searchterm)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	response, err := r.service.SearchVehicles(rs, searchterm, paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

// listViolations
func (r *vehicleResource) listRecentViolations(c *routing.Context) error {
	// get user id
	uid := c.Query("uid", "0")
	fmt.Println(uid)
	offset := 0
	limit := 50
	resp, err := r.service.ListRecentViolations(app.GetRequestScope(c), offset, limit, uid)
	if err != nil {
		return err
	}
	return c.Write(resp)
}

// getCurrentViolations.. single violation that has just happened
func (r *vehicleResource) getCurrentViolations(c *routing.Context) error {
	resp, err := r.service.GetCurrentViolations(app.GetRequestScope(c))
	if err != nil {
		return err
	}
	return c.Write(resp)
}

// getUnavailable
func (r *vehicleResource) getUnavailable(c *routing.Context) error {
	resp, err := r.service.GetUnavailableDevices(app.GetRequestScope(c))
	if err != nil {
		return err
	}
	return c.Write(resp)
}
