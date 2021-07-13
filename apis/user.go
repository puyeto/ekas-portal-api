package apis

import (
	"crypto/rand"
	"math/big"
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/errors"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

type (
	// userService specifies the interface for the user service needed by userResource.
	userService interface {
		// GetUser returns the user with the specified user ID.
		GetUser(rs app.RequestScope, id uint32) (*models.AuthUsers, error)
		Register(rs app.RequestScope, usr *models.AdminUserDetails) (int32, error)
		Login(rs app.RequestScope, usr *models.Credential) (*models.AdminUserDetails, error)
		SubmitUserRole(rs app.RequestScope, usr *models.AdminUserRoles) (*models.AdminUserRoles, error)
		Delete(rs app.RequestScope, id int32) error
		Query(rs app.RequestScope, offset, limit, cid int) ([]models.AuthUsers, error)
		Count(rs app.RequestScope, cid int) (int, error)
		Update(rs app.RequestScope, model *models.AuthUsers) (*models.AuthUsers, error)
		ResetPassword(rs app.RequestScope, model *models.ResetPassword) error
		QueryDepartments(rs app.RequestScope) ([]models.Departments, error)
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
	rg.Get("/users/list", r.query)
	rg.Post("/register", r.register)
	rg.Post("/login", r.login)
	rg.Put("/reset-password", r.resetPassword)
	rg.Put("/users/update", r.update)
	rg.Get("/ping", r.healthCheck)
	rg.Get("/otp-request", r.OTPRequest)
	rg.Get("/departments/list", r.queryDepartments)
}

func (r *userResource) query(c *routing.Context) error {
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

func (r *userResource) getuser(c *routing.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	response, err := r.service.GetUser(app.GetRequestScope(c), uint32(id))
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

// Reset admin password
func (r *userResource) resetPassword(c *routing.Context) error {
	var model models.ResetPassword
	if err := c.Read(&model); err != nil {
		return err
	}

	err := r.service.ResetPassword(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(map[string]int32{
		"last_insert_id": model.UserID,
	})
}

func (r *userResource) update(c *routing.Context) error {
	var model models.AuthUsers
	if err := c.Read(&model); err != nil {
		return err
	}

	response, err := r.service.Update(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *userResource) queryDepartments(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	items, err := r.service.QueryDepartments(rs)
	if err != nil {
		return err
	}
	return c.Write(items)
}

// OTPRequest ...
func (r *userResource) OTPRequest(c *routing.Context) error {
	phone := c.Query("phone", "0")
	if phone == "0" {
		return errors.Unauthorized("Invalid Phone Number")
	}
	otp, err := getRandNum()
	if err != nil {
		return errors.Unauthorized("OTP Error")
	}

	// send sms
	app.MessageChan <- app.MessageDetails{
		Message:  strconv.Itoa(otp),
		ToNumber: phone,
		Type:     "OTPRequest",
	}

	return c.Write(map[string]int{
		"otp": otp,
	})
}

func getRandNum() (int, error) {
	nBig, e := rand.Int(rand.Reader, big.NewInt(8999))
	if e != nil {
		return 0, e
	}
	return int(nBig.Int64() + 1000), nil
}

func (r *userResource) healthCheck(c *routing.Context) error {
	return c.Write(map[string]interface{}{
		"status":  200,
		"message": "Health check successfull",
	})
}
