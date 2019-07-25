package services

import (
	"errors"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// settingDAO specifies the interface of the setting DAO needed by SettingService.
type settingDAO interface {
	// Get returns the setting with the specified setting ID.
	Get(rs app.RequestScope, id int) (*models.Settings, error)
	// Count returns the number of settings.
	Count(rs app.RequestScope) (int, error)
	// Query returns the list of settings with the given offset and limit.
	Query(rs app.RequestScope, offset, limit int) ([]models.Settings, error)
	// Create saves a new setting in the storage.
	Create(rs app.RequestScope, setting *models.Settings) error
	// Update updates the setting with given ID in the storage.
	Update(rs app.RequestScope, id int, setting *models.Settings) error
	// Delete removes the setting with given ID from the storage.
	Delete(rs app.RequestScope, id int) error
	GenerateKey(rs app.RequestScope, keys []string, assignto int) error
	CountKeys(rs app.RequestScope) (int, error)
	QueryKeys(rs app.RequestScope, offset, limit int) ([]models.LicenseKeys, error)
	AssignKey(rs app.RequestScope, model *models.LicenseKeys) error
	GetKey(rs app.RequestScope, key string) (*models.LicenseKeys, error)
	UpdateKey(rs app.RequestScope, model *models.LicenseKeys) error
}

// SettingService provides services related with settings.
type SettingService struct {
	dao settingDAO
}

// NewSettingService creates a new SettingService with the given setting DAO.
func NewSettingService(dao settingDAO) *SettingService {
	return &SettingService{dao}
}

// Get returns the setting with the specified the setting ID.
func (s *SettingService) Get(rs app.RequestScope, id int) (*models.Settings, error) {
	return s.dao.Get(rs, id)
}

// Create creates a new setting.
func (s *SettingService) Create(rs app.RequestScope, model *models.Settings) (*models.Settings, error) {
	if err := model.ValidateSettings(); err != nil {
		return nil, err
	}
	if err := s.dao.Create(rs, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, model.SettingID)
}

// Update updates the setting with the specified ID.
func (s *SettingService) Update(rs app.RequestScope, id int, model *models.Settings) (*models.Settings, error) {
	if err := model.ValidateSettings(); err != nil {
		return nil, err
	}
	if err := s.dao.Update(rs, id, model); err != nil {
		return nil, err
	}
	return s.dao.Get(rs, id)
}

// Delete deletes the setting with the specified ID.
func (s *SettingService) Delete(rs app.RequestScope, id int) (*models.Settings, error) {
	setting, err := s.dao.Get(rs, id)
	if err != nil {
		return nil, err
	}
	err = s.dao.Delete(rs, id)
	return setting, err
}

// Count returns the number of settings.
func (s *SettingService) Count(rs app.RequestScope) (int, error) {
	return s.dao.Count(rs)
}

// Query returns the settings with the specified offset and limit.
func (s *SettingService) Query(rs app.RequestScope, offset, limit int) ([]models.Settings, error) {
	return s.dao.Query(rs, offset, limit)
}

// GenerateKey save generated keys
func (s *SettingService) GenerateKey(rs app.RequestScope, model *models.GenKeys) ([]string, error) {
	if err := model.ValidateGenKeys(); err != nil {
		return nil, err
	}

	if model.AssignTo == 0 {
		return nil, errors.New("Assign To is required")
	}

	var letterBytes = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	if model.Type == "NUMERIC" {
		letterBytes = "0123456789"
	} else if model.Type == "ALPHABETIC" {
		letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	var keys []string
	for i := 0; i < model.Number; i++ {
		str := app.RandStringBytes(letterBytes, 4) + "-" + app.RandStringBytes(letterBytes, 4) + "-" + app.RandStringBytes(letterBytes, 4) + "-" + app.RandStringBytes(letterBytes, 4)
		keys = append(keys, str)
	}

	if len(keys) > 0 {
		if err := s.dao.GenerateKey(rs, keys, model.AssignTo); err != nil {
			return nil, err
		}
	}

	return keys, nil
}

// CountKeys returns the number of keys.
func (s *SettingService) CountKeys(rs app.RequestScope) (int, error) {
	return s.dao.CountKeys(rs)
}

// QueryKeys return keys with the specified offset and limit.
func (s *SettingService) QueryKeys(rs app.RequestScope, offset, limit int) ([]models.LicenseKeys, error) {
	return s.dao.QueryKeys(rs, offset, limit)
}

// AssignKey assign generated keys
func (s *SettingService) AssignKey(rs app.RequestScope, model *models.LicenseKeys) error {
	if err := model.ValidateLicenseKeys(); err != nil {
		return err
	}

	if err := s.dao.AssignKey(rs, model); err != nil {
		return err
	}

	return nil
}

// GetKey returns the keys with the specified the key string.
func (s *SettingService) GetKey(rs app.RequestScope, key string) (*models.LicenseKeys, error) {
	return s.dao.GetKey(rs, key)
}

// UpdateKey updates the keys.
func (s *SettingService) UpdateKey(rs app.RequestScope, model *models.LicenseKeys) (*models.LicenseKeys, error) {
	if err := model.ValidateLicenseKeys(); err != nil {
		return nil, err
	}
	if err := s.dao.UpdateKey(rs, model); err != nil {
		return nil, err
	}
	return s.dao.GetKey(rs, model.KeyString)
}
