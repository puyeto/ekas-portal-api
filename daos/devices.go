package daos

import (
	"strconv"
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
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
	err := rs.Tx().Select().From("device_details").Model(id, &device)
	return &device, err
}

// Create saves a new device record in the database.
// The Device.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *DeviceDAO) Create(rs app.RequestScope, device *models.Devices) error {
	res, err := rs.Tx().Insert("device_details", dbx.Params{
		"device_id":           device.DeviceID,
		"device_name":         device.DeviceName,
		"device_serial_no":    strings.ToUpper(device.DeviceSerialNo),
		"device_model":        strings.ToUpper(device.DeviceModelNo),
		"device_manufacturer": strings.ToUpper(device.DeviceManufacturer),
		"configured":          device.Configured,
		"note":                device.Note}).Execute()
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	device.ID = int32(id)
	return err
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
func (dao *DeviceDAO) CountConfiguredDevices(rs app.RequestScope, vehicleid, deviceid int) (int, error) {
	var count int
	query := "SELECT COUNT(device_id) FROM vehicle_configuration"
	if vehicleid > 0 && deviceid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "' AND device_id = '" + strconv.Itoa(deviceid) + "'"
	} else if vehicleid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "'"
	} else if deviceid > 0 {
		query += " WHERE device_id = '" + strconv.Itoa(deviceid) + "'"
	}
	q := rs.Tx().NewQuery(query)
	err := q.Row(&count)

	return count, err
}

// ConfiguredDevices retrieves the device records with the specified offset and limit from the database.
func (dao *DeviceDAO) ConfiguredDevices(rs app.RequestScope, offset, limit, vehicleid, deviceid int) ([]models.DeviceConfiguration, error) {
	devices := []models.DeviceConfiguration{}
	query := "SELECT device_id, vehicle_id, data->'$.sim_imei' AS sim_imei, created_on FROM vehicle_configuration"
	if vehicleid > 0 && deviceid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "' AND device_id = '" + strconv.Itoa(deviceid) + "'"
	} else if vehicleid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "'"
	} else if deviceid > 0 {
		query += " WHERE device_id = '" + strconv.Itoa(deviceid) + "'"
	}
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit)
	q := rs.Tx().NewQuery(query)
	err := q.All(&devices)

	return devices, err
}
