package apis

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// vehicleService specifies the interface for the vehicle service needed by vehicleResource.
	vehicleService interface {
		Create(rs app.RequestScope, model *models.Vehicle) (int, error)
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
