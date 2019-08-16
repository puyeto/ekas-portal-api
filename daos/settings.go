package daos

import (
	"strconv"
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/errors"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// SettingDAO persists setting data in database
type SettingDAO struct{}

// NewSettingDAO creates a new SettingDAO
func NewSettingDAO() *SettingDAO {
	return &SettingDAO{}
}

// Get reads the setting with the specified ID from the database.
func (dao *SettingDAO) Get(rs app.RequestScope, id int) (*models.Settings, error) {
	var setting models.Settings
	err := rs.Tx().Select().Model(id, &setting)
	return &setting, err
}

// Update saves the changes to an setting in the database.
func (dao *SettingDAO) Update(rs app.RequestScope, id int, setting *models.Settings) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	setting.SettingID = id
	return rs.Tx().Model(setting).Exclude("Id").Update()
}

// GenerateKey ...
func (dao *SettingDAO) GenerateKey(rs app.RequestScope, keys []string, assignto int) error {
	var query = "INSERT INTO license_keys(key_string, assign_to) VALUES"
	for _, key := range keys {
		query += "('" + key + "', '" + strconv.Itoa(assignto) + "'),"
	}

	//trim the last ,
	query = strings.TrimSuffix(query, ",")

	q := rs.Tx().NewQuery(query)
	_, err := q.Execute()

	return err
}

// CountKeys returns the number of key records in the database.
func (dao *SettingDAO) CountKeys(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("license_keys").Row(&count)
	return count, err
}

// QueryKeys retrieves the keys records with the specified offset and limit from the database.
func (dao *SettingDAO) QueryKeys(rs app.RequestScope, offset, limit int) ([]models.LicenseKeys, error) {
	keys := []models.LicenseKeys{}
	err := rs.Tx().Select("key_string", "assign_to", "license_keys.status", "COALESCE(company_id, 0) AS company_id", "COALESCE(company_name, '') AS company_name", "COALESCE(company_contacts, '') AS company_contacts").
		// LeftJoin("company_users", dbx.NewExp("company_users.user_id = license_keys.assign_to")).
		LeftJoin("companies", dbx.NewExp("companies.company_id = license_keys.assign_to")).
		Offset(int64(offset)).Limit(int64(limit)).All(&keys)
	return keys, err
}

// GetKey reads the keys with the specified ID from the database.
func (dao *SettingDAO) GetKey(rs app.RequestScope, key string) (*models.LicenseKeys, error) {
	var model models.LicenseKeys
	err := rs.Tx().Select().Model(key, &model)

	return &model, err
}

// AssignKey assign key to user by user id.
func (dao *SettingDAO) AssignKey(rs app.RequestScope, model *models.LicenseKeys) error {
	keydata, err := dao.GetKey(rs, model.KeyString)
	if err != nil {
		return err
	}

	// check if a key has been assigned
	if keydata.Status > 0 {
		return errors.New("Key Already Assigned")
	}

	_, err = rs.Tx().Update("license_keys", dbx.Params{
		"assign_to": model.AssignTo,
		"status":    1,
	}, dbx.HashExp{"key_string": model.KeyString}).Execute()
	if err != nil {
		return err
	}

	return nil
}

// UpdateKey saves the changes to keys in the database.
func (dao *SettingDAO) UpdateKey(rs app.RequestScope, model *models.LicenseKeys) error {
	return rs.Tx().Model(model).Exclude("KeyString", "CompanyID", "CompanyName", "CompanyContacts").Update()
}
