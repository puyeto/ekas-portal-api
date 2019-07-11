package apis

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/errors"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// userService specifies the interface for the user service needed by userResource.
	userService interface {
		// GetUser returns the user with the specified user ID.
		GetUser(rs app.RequestScope, id int32) (*models.ListUserDetails, error)
		Register(rs app.RequestScope, usr *models.AdminUserDetails) (int32, error)
		Login(rs app.RequestScope, usr *models.Credential) (*models.AdminUserDetails, error)
		SubmitUserRole(rs app.RequestScope, usr *models.AdminUserRoles) (*models.AdminUserRoles, error)
		Delete(rs app.RequestScope, id int32) error
	}

	// userResource defines the handlers for the CRUD APIs.
	userResource struct {
		service userService
	}
)

// ServeUserResource sets up the routing of user endpoints and the corresponding handlers.
func ServeUserResource(rg *routing.RouteGroup, service userService) {
	r := &userResource{service}
	rg.Get("/user/<id>", r.getuser)
	rg.Post("/user/role", r.submitroles)
	rg.Delete("/user/delete/<id>", r.delete)
	rg.Post("/register", r.register)
	rg.Post("/login", r.login)
}

func (r *userResource) getuser(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.GetUser(app.GetRequestScope(c), int32(id))
	if err != nil {
		return err
	}

	return c.Write(response)
}

// login -user login
func (r *userResource) login(c *routing.Context) error {
	var credential models.Credential
	if err := c.Read(&credential); err != nil {
		return errors.BadRequest(err.Error())
	}

	identity, err := r.service.Login(app.GetRequestScope(c), &credential)
	if err != nil {
		return errors.Unauthorized(err.Error())
	}

	return c.Write(identity)

}

func (r *userResource) register(c *routing.Context) error {
	var model models.AdminUserDetails
	if err := c.Read(&model); err != nil {
		return err
	}

	response, err := r.service.Register(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(map[string]int32{
		"last_insert_id": response,
	})
}

func (r *userResource) submitroles(c *routing.Context) error {
	var model models.AdminUserRoles
	if err := c.Read(&model); err != nil {
		return err
	}
	response, err := r.service.SubmitUserRole(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *userResource) delete(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	if err := r.service.Delete(app.GetRequestScope(c), int32(id)); err != nil {
		return err
	}

	return c.Write(map[string]string{
		"message": "Record Deleted Successfully",
	})
}
