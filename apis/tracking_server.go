package apis

import (
	"fmt"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type (
	// trackingServerService specifies the interface for the trackingServer service needed by trackingServerResource.
	trackingServerService interface {
		TrackingServerLogin(rs app.RequestScope, model *models.TrackingServerAuth) (models.AdminUserDetails, error)
		TrackingServerUserDevices(rs app.RequestScope, model *models.UserData) ([]models.VehicleDetails, error)
		TrackingServerAddDevices(rs app.RequestScope, model *models.AddDeviceDetails, lang string, userhash string) (interface{}, error)
		TrackingServerEditDevices(rs app.RequestScope, model *models.AddDeviceDetails, lang string, userhash string) (interface{}, error)
	}

	// trackingServerResource defines the handlers for the CRUD APIs.
	trackingServerResource struct {
		service trackingServerService
	}
)

// ServeTrackingServerResource sets up the routing of trackingServer endpoints and the corresponding handlers.
func ServeTrackingServerResource(rg *routing.RouteGroup, service trackingServerService) {
	r := &trackingServerResource{service}
	rg.Post("/trackingserverlogin", r.trackingServerLogin)
	rg.Post("/trackingservergetdevices", r.trackingservergetuserdevices)
	rg.Post("/trackingserveradddevices", r.trackingserveradduserdevices)
	rg.Put("/trackingservereditdevices", r.trackingserveredituserdevices)
}

func (r *trackingServerResource) trackingServerLogin(c *routing.Context) error {
	var model models.TrackingServerAuth
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.TrackingServerLogin(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *trackingServerResource) trackingservergetuserdevices(c *routing.Context) error {
	var model models.UserData
	if err := c.Read(&model); err != nil {
		return err
	}

	response, err := r.service.TrackingServerUserDevices(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

// trackingserveradduserdevices ...
func (r *trackingServerResource) trackingserveradduserdevices(c *routing.Context) error {
	var model models.AddDevice
	if err := c.Read(&model); err != nil {
		return err
	}

	fmt.Println(model)

	// valid device structurs
	var u = model.UserData
	if err := u.ValidateUserData(); err != nil {
		return err
	}

	// valid device structurs
	var v = model.DeviceData
	if err := v.ValidateAddDevices(); err != nil {
		return err
	}

	response, err := r.service.TrackingServerAddDevices(app.GetRequestScope(c), &model.DeviceData, u.Lang, u.UserHash)
	if err != nil {
		return err
	}

	return c.Write(response)
}

// trackingserveredituserdevices ...
func (r *trackingServerResource) trackingserveredituserdevices(c *routing.Context) error {
	var model models.AddDevice
	if err := c.Read(&model); err != nil {
		return err
	}

	// valid device structurs
	var u = model.UserData
	if err := u.ValidateUserData(); err != nil {
		return err
	}

	// valid device structurs
	var v = model.DeviceData
	if err := v.ValidateAddDevices(); err != nil {
		return err
	}

	response, err := r.service.TrackingServerEditDevices(app.GetRequestScope(c), &model.DeviceData, u.Lang, u.UserHash)
	if err != nil {
		return err
	}

	return c.Write(response)
}
