package apis

import (
	"net/http"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	routing "github.com/go-ozzo/ozzo-routing"
)

type (
	// settingService specifies the interface for the setting service needed by settingResource.
	settingService interface {
		Get(rs app.RequestScope, id int) (*models.Settings, error)
		Update(rs app.RequestScope, id int, model *models.Settings) (*models.Settings, error)
		GenerateKey(rs app.RequestScope, model *models.GenKeys) ([]string, error)
		CountKeys(rs app.RequestScope) (int, error)
		QueryKeys(rs app.RequestScope, offset, limit int) ([]models.LicenseKeys, error)
		AssignKey(rs app.RequestScope, model *models.LicenseKeys) error
		GetKey(rs app.RequestScope, key string) (*models.LicenseKeys, error)
		UpdateKey(rs app.RequestScope, model *models.LicenseKeys) (*models.LicenseKeys, error)
	}

	// settingResource defines the handlers for the CRUD APIs.
	settingResource struct {
		service settingService
	}
)

// ServeSettingResource sets up the routing of setting endpoints and the corresponding handlers.
func ServeSettingResource(rg *routing.RouteGroup, service settingService) {
	r := &settingResource{service}
	rg.Get("/settings/get", r.get)
	rg.Put("/settings/update", r.update)
	rg.Post("/settings/generate-keys", r.generateKey)
	rg.Put("/settings/keys/update", r.updateKey)
	rg.Get("/setting/list-keys", r.queryKeys)
	rg.Post("/setting/assign-key", r.assignKey)
}

func (r *settingResource) get(c *routing.Context) error {
	response, err := r.service.Get(app.GetRequestScope(c), 1)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *settingResource) update(c *routing.Context) error {
	var id = 1

	rs := app.GetRequestScope(c)

	model, err := r.service.Get(rs, id)
	if err != nil {
		return err
	}

	if err := c.Read(model); err != nil {
		return err
	}

	response, err := r.service.Update(rs, id, model)
	if err != nil {
		return err
	}

	return c.Write(response)
}

func (r *settingResource) generateKey(c *routing.Context) error {
	var model models.GenKeys
	if err := c.Read(&model); err != nil {
		return err
	}

	keys, err := r.service.GenerateKey(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(keys)
}

func (r *settingResource) queryKeys(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	count, err := r.service.CountKeys(rs)
	if err != nil {
		return err
	}
	paginatedList := getPaginatedListFromRequest(c, count)
	items, err := r.service.QueryKeys(app.GetRequestScope(c), paginatedList.Offset(), paginatedList.Limit())
	if err != nil {
		return err
	}
	paginatedList.Items = items
	return c.Write(paginatedList)
}

func (r *settingResource) assignKey(c *routing.Context) error {
	var model models.LicenseKeys
	if err := c.Read(&model); err != nil {
		return err
	}

	err := r.service.AssignKey(app.GetRequestScope(c), &model)
	if err != nil {
		return err
	}

	return c.Write(http.StatusOK)
}

func (r *settingResource) updateKey(c *routing.Context) error {
	rs := app.GetRequestScope(c)
	var model models.LicenseKeys
	if err := c.Read(&model); err != nil {
		return err
	}

	_, err := r.service.GetKey(rs, model.KeyString)
	if err != nil {
		return err
	}

	response, err := r.service.UpdateKey(rs, &model)
	if err != nil {
		return err
	}

	return c.Write(response)
}
