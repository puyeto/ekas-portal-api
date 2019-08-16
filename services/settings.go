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
	// Update updates the setting with given ID in the storage.
	Update(rs app.RequestScope, id int, setting *models.Settings) error
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
