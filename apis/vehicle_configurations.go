package apis

import (
	"strconv"
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type (
	// vehicleService specifies the interface for the vehicle service needed by vehicleResource.
	vehicleService interface {
		GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error)
		GetConfigurationDetails(rs app.RequestScope, vehicleid, deviceid int) (*models.VehicleConfigDetails, error)
		GetTripDataByDeviceID(deviceid string, offset, limit int, orderby string) ([]models.DeviceData, error)
		GetTripDataByDeviceIDBtwDates(deviceid string, offset, limit int, from, to int64) ([]models.DeviceData, error)
		Create(rs app.RequestScope, model *models.Vehicle) (int, error)
		CountTripRecords(rs app.RequestScope, deviceid string) (int, error)
		CountRedisTripRecords(deviceid string) int

		CountOverspeed(rs app.RequestScope, deviceid string) (int, error)
		CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error)
		GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.DeviceData, error)
		GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.DeviceData, error)
		GetCurrentViolations(rs app.RequestScope) (models.DeviceData, error)
		ListAllViolations(rs app.RequestScope, offset, limit int) ([]models.CurrentViolations, error)
		XMLListAllViolations(rs app.RequestScope, offset, limit int) ([]models.XMLResults, error)
		CountAllViolations() (int, error)
		SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int, qtype string) ([]models.SearchDetails, error)
		CountSearches(rs app.RequestScope, searchterm, qtype string) (int, error)
		GetUnavailableDevices(rs app.RequestScope) ([]models.DeviceData, error)
		GetOfflineViolations(rs app.RequestScope, deviceid string) ([]models.DeviceData, error)
		CountTripDataByDeviceID(deviceid string) (int, error)
		CountTripRecordsBtwDates(deviceid string, from int64, to int64) (int, error)
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
	rg.Get("/device/data-range/<id>", r.getTripDataByDeviceIDBtwDates)
	rg.Get("/getoverspeed/<id>", r.getOverspeedsByDeviceID)
	rg.Get("/getfailsafe/<id>", r.getFailsafeByDeviceID)
	rg.Get("/getdisconnects/<id>", r.getDisconnectsByDeviceID)
	rg.Get("/currentviolation", r.getCurrentViolations)
	rg.Get("/listviolations", r.listAllViolations)
	rg.Get("/xmllistviolations", r.xmlListRecentViolations)
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
	orderby := c.Query("orderby", "desc")

	// rs := app.GetRequestScope(c)
	count, _ := r.service.CountTripDataByDeviceID(deviceid)
	paginatedList := getPaginatedListFromRequest(c, count)

	response, err := r.service.GetTripDataByDeviceID(deviceid, paginatedList.Offset(), paginatedList.Limit(), orderby)
	if err != nil {
		return err
	}

	paginatedList.Items = response
	return c.Write(paginatedList)
}

// getTripDataByDeviceID ...
func (r *vehicleResource) getTripDataByDeviceIDBtwDates(c *routing.Context) error {
	var model models.TripBetweenDates
	model.DeviceID = c.Param("id")
	start, _ := strconv.Atoi(c.Query("start", "0"))
	model.From = int64(start)
	stop, _ := strconv.Atoi(c.Query("stop", "0"))
	model.To = int64(stop)

	count, _ := r.service.CountTripRecordsBtwDates(model.DeviceID, model.From, model.To)
	paginatedList := getPaginatedListFromRequest(c, count)

	response, err := r.service.GetTripDataByDeviceIDBtwDates(model.DeviceID, paginatedList.Offset(), paginatedList.Limit(), model.From, model.To)
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

// listAllViolations ...
func (r *vehicleResource) listAllViolations(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	count, _ := r.service.CountAllViolations()
	paginatedList := getPaginatedListFromRequest(c, count)

	response, err := r.service.ListAllViolations(rs, paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
}

// xmlListRecentViolations
func (r *vehicleResource) xmlListRecentViolations(c *routing.Context) error {

	rs := app.GetRequestScope(c)
	count, _ := r.service.CountAllViolations()
	paginatedList := getPaginatedListFromRequest(c, count)

	response, err := r.service.XMLListAllViolations(rs, paginatedList.Offset(), paginatedList.Limit())
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
	if count > 0 {
		response, err := r.service.GetOverspeedByDeviceID(rs, deviceid, paginatedList.Offset(), paginatedList.Limit())
		if err != nil {
			return err
		}
		paginatedList.Items = response
	}
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
	qtype := c.Query("type", "")
	rs := app.GetRequestScope(c)
	count, err := r.service.CountSearches(rs, searchterm, qtype)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	response, err := r.service.SearchVehicles(rs, searchterm, paginatedList.Offset(), paginatedList.Limit(), qtype)
	if err != nil {
		return err
	}
	paginatedList.Items = response
	return c.Write(paginatedList)
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
