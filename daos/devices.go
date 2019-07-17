package daos

import (
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// DeviceDAO persists device data in database
type DeviceDAO struct{}

// NewDeviceDAO creates a new DeviceDAO
func NewDeviceDAO() *DeviceDAO {
	return &DeviceDAO{}
}

// Get reads the device with the specified ID from the database.
func (dao *DeviceDAO) Get(rs app.RequestScope, id int32) (*models.Devices, error) {
	var device models.Devices
	err := rs.Tx().Select().Model(id, &device)
	return &device, err
}

// Create saves a new device record in the database.
// The Device.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *DeviceDAO) Create(rs app.RequestScope, device *models.Devices) error {
	device.DeviceID = 0
	return rs.Tx().Model(device).Insert()
}

// Update saves the changes to an device in the database.
func (dao *DeviceDAO) Update(rs app.RequestScope, id int32, device *models.Devices) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	device.DeviceID = id
	return rs.Tx().Model(device).Exclude("Id").Update()
}

// Delete deletes an device with the specified ID from the database.
func (dao *DeviceDAO) Delete(rs app.RequestScope, id int32) error {
	device, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(device).Delete()
}

// Count returns the number of the device records in the database.
func (dao *DeviceDAO) Count(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("device_details").Row(&count)
	return count, err
}

// Query retrieves the device records with the specified offset and limit from the database.
func (dao *DeviceDAO) Query(rs app.RequestScope, offset, limit int) ([]models.Devices, error) {
	devices := []models.Devices{}
	err := rs.Tx().Select().From("device_details").OrderBy("id").Offset(int64(offset)).Limit(int64(limit)).All(&devices)
	return devices, err
}

// CountConfiguredDevices returns the number of the device configuration records in the database.
func (dao *DeviceDAO) CountConfiguredDevices(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(distinct device_id)").From("vehicle_configuration").Row(&count)
	return count, err
}

// Query retrieves the device records with the specified offset and limit from the database.
func (dao *DeviceDAO) ConfiguredDevices(rs app.RequestScope, offset, limit int) ([]models.DeviceConfiguration, error) {
	devices := []models.DeviceConfiguration{}
	query := "SELECT DISTINCT device_id, data->'$.sim_imei' AS sim_imei FROM vehicle_configuration"
	query += " LIMIT 0, 100"
	q := rs.Tx().NewQuery(query)
	err := q.All(&devices)

	return devices, err
}
