package apis

import (
	"fmt"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// trackingServerService specifies the interface for the trackingServer service needed by trackingServerResource.
	trackingServerService interface {
		TrackingServerLogin(rs app.RequestScope, model *models.TrackingServerAuth) (interface{}, error)
		TrackingServerUserDevices(rs app.RequestScope, model *models.UserData) (interface{}, error)
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

	fmt.Println(model.Lang, model.UserHash)

	response, err := r.service.TrackingServerUserDevices(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}
