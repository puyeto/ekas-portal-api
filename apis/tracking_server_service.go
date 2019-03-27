package apis

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// trackingServerServiceService specifies the interface for the trackingServerService service needed by trackingServerServiceResource.
	trackingServerServiceService interface {
		TrackingServerGetServices(rs app.RequestScope, model *models.UserData) (interface{}, error)
		TrackingServerAddServices(rs app.RequestScope, model *models.AddTrackingServiceDetails, lang string, userhash string, deviceid string) (interface{}, error)
	}

	// trackingServerServiceResource defines the handlers for the CRUD APIs.
	trackingServerServiceResource struct {
		service trackingServerServiceService
	}
)

// ServeTrackingServerService sets up the routing of trackingServerService endpoints and the corresponding handlers.
func ServeTrackingServerServiceResource(rg *routing.RouteGroup, service trackingServerServiceService) {
	r := &trackingServerServiceResource{service}
	rg.Post("/trackingservergetservices", r.get)
	rg.Post("/trackingserveraddservices", r.add)
}

func (r *trackingServerServiceResource) get(c *routing.Context) error {
	var model models.UserData
	if err := c.Read(&model); err != nil {
		return err
	}

	response, err := r.service.TrackingServerGetServices(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *trackingServerServiceResource) add(c *routing.Context) error {
	var model models.AddServices
	if err := c.Read(&model); err != nil {
		return err
	}

	// valid device structurs
	var u = model.UserData
	if err := u.ValidateUserData(); err != nil {
		return err
	}

	// valid device structurs
	var v = model.AddTrackingServiceDetails
	if err := v.ValidateAddTrackingServiceDetails(); err != nil {
		return err
	}

	response, err := r.service.TrackingServerAddServices(app.GetRequestScope(c), &v, u.Lang, u.UserHash, string(v.DeviceID))
	if err != nil {
		return err
	}

	return c.Write(response)
}
