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
	err := rs.Tx().Select().From("device_details").OrderBy("device_id ASC").Offset(int64(offset)).Limit(int64(limit)).All(&devices)
	return devices, err
}

// CountConfiguredDevices returns the number of the device configuration records in the database.
func (dao *DeviceDAO) CountConfiguredDevices(rs app.RequestScope, vehicleid, deviceid int) (int, error) {
	var count int
	query := "SELECT COUNT(device_id) FROM vehicle_configuration"
	if vehicleid > 0 && deviceid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "' AND device_id = '" + strconv.Itoa(deviceid) + "' AND status=1"
	} else if vehicleid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "' AND status=1"
	} else if deviceid > 0 {
		query += " WHERE device_id = '" + strconv.Itoa(deviceid) + "' AND status=1"
	} else {
		query += " WHERE status=1 "
	}
	q := rs.Tx().NewQuery(query)
	err := q.Row(&count)

	return count, err
}

// ConfiguredDevices retrieves the device records with the specified offset and limit from the database.
func (dao *DeviceDAO) ConfiguredDevices(rs app.RequestScope, offset, limit, vehicleid, deviceid int) ([]models.DeviceConfiguration, error) {
	devices := []models.DeviceConfiguration{}
	query := "SELECT vc.device_id, vehicle_id, COALESCE(JSON_VALUE(data,'$.device_detail.registration_no')) AS device_name, JSON_VALUE(data, '$.sim_imei') AS sim_imei, vc.created_on, vc.status AS status, "
	query += " JSON_VALUE(data, '$.device_detail.chasis_no') AS chassis_no, JSON_VALUE(data, '$.device_detail.make_type') AS make_type, JSON_VALUE(data, '$.device_detail.device_type') AS device_type, "
	query += " JSON_VALUE(data, '$.device_detail.serial_no') AS serial_no, JSON_VALUE(data, '$.device_detail.preset_speed') AS preset_speed, JSON_VALUE(data, '$.device_detail.set_frequency') AS set_frequency, "
	query += " JSON_VALUE(data, '$.device_detail.fitting_date') AS fitting_date, DATE_ADD(JSON_VALUE(data, '$.device_detail.fitting_date'), INTERVAL 1 YEAR) AS expiry_date, JSON_VALUE(data, '$.device_detail.fitting_center') AS fitting_center, "
	query += " JSON_VALUE(data, '$.device_detail.certificate') AS certificate, JSON_VALUE(data, '$.device_detail.email_address') AS email_address, JSON_VALUE(data, '$.device_detail.agent_phone') AS agent_phone, "
	query += " JSON_VALUE(data, '$.device_detail.agent_location') AS agent_location, JSON_VALUE(data, '$.device_detail.owner_name') AS owner_name, JSON_VALUE(data, '$.device_detail.owner_phone_number') AS owner_phone_number,"
	query += " COALESCE(dd.status, 0) AS device_status, COALESCE(dd.status_reason, 'Device record does not exit') AS reason FROM vehicle_configuration AS vc"
	query += " LEFT JOIN device_details AS dd ON (dd.device_id = vc.device_id)"

	if vehicleid > 0 && deviceid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "' AND device_id = '" + strconv.Itoa(deviceid) + "' AND vc.status=1"
	} else if vehicleid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "' AND vc.status=1"
	} else if deviceid > 0 {
		query += " WHERE device_id = '" + strconv.Itoa(deviceid) + "' AND vc.status=1"
	} else {
		query += " WHERE vc.status=1 "
	}
	query += " LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit)
	q := rs.Tx().NewQuery(query)
	err := q.All(&devices)

	return devices, err
}
