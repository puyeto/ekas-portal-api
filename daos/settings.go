package daos

import (
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
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

// Create saves a new setting record in the database.
// The Setting.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *SettingDAO) Create(rs app.RequestScope, setting *models.Settings) error {
	setting.SettingID = 0
	return rs.Tx().Model(setting).Insert()
}

// Update saves the changes to an setting in the database.
func (dao *SettingDAO) Update(rs app.RequestScope, id int, setting *models.Settings) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	setting.SettingID = id
	return rs.Tx().Model(setting).Exclude("Id").Update()
}

// Delete deletes an setting with the specified ID from the database.
func (dao *SettingDAO) Delete(rs app.RequestScope, id int) error {
	setting, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(setting).Delete()
}

// Count returns the number of the setting records in the database.
func (dao *SettingDAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("setting").Row(&count)
	return count, err
}

// Query retrieves the setting records with the specified offset and limit from the database.
func (dao *SettingDAO) Query(rs app.RequestScope, offset, limit int) ([]models.Settings, error) {
	settings := []models.Settings{}
	err := rs.Tx().Select().OrderBy("id").Offset(int64(offset)).Limit(int64(limit)).All(&settings)
	return settings, err
}

// GenerateKey ...
func (dao *SettingDAO) GenerateKey(rs app.RequestScope, keys []string) error {
	var query = "INSERT INTO license_keys(key_string) VALUES"
	for _, key := range keys {
		query += "('" + key + "'),"
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
	err := rs.Tx().Select().Offset(int64(offset)).Limit(int64(limit)).All(&keys)
	return keys, err
}
