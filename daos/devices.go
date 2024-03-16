package daos

import (
	"encoding/json"
	"fmt"
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
		"sim_serial_no":       device.SimSerialNo,
		"certificate_no":      device.CertificateNo,
		"sim_number":          device.SimNumber,
		"motherboard_no":      device.MotherboardNO,
		"technician":          device.Technician,
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
	device.ID = id
	// return rs.Tx().Model(device).Exclude("ID", "CreatedOn", "Positions").Update()
	_, err := rs.Tx().Update("device_details", dbx.Params{
		"device_id":           device.DeviceID,
		"device_name":         device.DeviceName,
		"device_serial_no":    strings.ToUpper(device.DeviceSerialNo),
		"device_model":        strings.ToUpper(device.DeviceModelNo),
		"device_manufacturer": strings.ToUpper(device.DeviceManufacturer),
		"sim_serial_no":       device.SimSerialNo,
		"certificate_no":      device.CertificateNo,
		"sim_number":          device.SimNumber,
		"motherboard_no":      device.MotherboardNO,
		"technician":          device.Technician,
		"configured":          device.Configured,
		"note":                device.Note},
		dbx.HashExp{"id": id}).Execute()
	return err

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
func (dao *DeviceDAO) Count(rs app.RequestScope, cid int) (int, error) {
	var count int
	var err error

	if cid == 0 {
		err = rs.Tx().Select("COUNT(*)").From("device_details").Row(&count)
	} else {
		err = rs.Tx().Select("COUNT(*)").From("device_details").Where(dbx.HashExp{"company_id": cid}).Row(&count)
	}

	return count, err
}

// Query retrieves the device records with the specified offset and limit from the database.
func (dao *DeviceDAO) Query(rs app.RequestScope, offset, limit, cid int) ([]models.Devices, error) {
	devices := []models.Devices{}
	var err error

	if cid == 0 {
		err = rs.Tx().Select("id", "device_id", "vehicle_id", "device_name", "device_details.company_id",
			"COALESCE(company_name, '') AS company_name", "device_serial_no", "device_model", "device_manufacturer",
			"sim_serial_no", "certificate_no", "sim_number", "motherboard_no", "technician", "status", "configured", "configuration_date", "status_reason", "created_on").
			LeftJoin("companies", dbx.NewExp("companies.company_id = device_details.company_id")).
			From("device_details").OrderBy("id DESC").Offset(int64(offset)).Limit(int64(limit)).All(&devices)

	} else {
		err = rs.Tx().Select("id", "device_id", "vehicle_id", "device_name", "device_details.company_id",
			"COALESCE(company_name, '') AS company_name", "device_serial_no", "device_model", "device_manufacturer",
			"sim_serial_no", "certificate_no", "sim_number", "motherboard_no", "technician", "status", "configured", "configuration_date", "status_reason", "created_on").
			LeftJoin("companies", dbx.NewExp("companies.company_id = device_details.company_id")).
			From("device_details").Where(dbx.HashExp{"device_details.company_id": cid}).
			OrderBy("id DESC").Offset(int64(offset)).Limit(int64(limit)).All(&devices)
	}

	return devices, err
}

// CountQueryPositions returns the number of the device position records in the database.
func (dao *DeviceDAO) CountQueryPositions(rs app.RequestScope, uid uint32) (int, error) {
	var count int
	var err error
	if uid > 0 {
		err = rs.Tx().Select("COUNT(device_id)").From("device_details").
			LeftJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = device_details.vehicle_id")).
			Where(dbx.NewExp("device_details.vehicle_id > 0")).
			Where(dbx.HashExp{"vehicle_details.user_id": uid}).Row(&count)
	} else {
		err = rs.Tx().Select("COUNT(device_id)").From("device_details").
			LeftJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = device_details.vehicle_id")).
			Where(dbx.NewExp("device_details.vehicle_id > 0")).Row(&count)
	}
	return count, err
}

// QueryPositions retrieves the device positions records with the specified offset and limit from the database.
func (dao *DeviceDAO) QueryPositions(rs app.RequestScope, offset, limit int, uid uint32, start, stop int64) ([]models.Devices, error) {
	var devices []models.Devices
	var d models.Devices
	var q *dbx.SelectQuery
	if uid > 0 {
		q = rs.Tx().Select("id", "device_id", "device_details.vehicle_id", "vehicle_reg_no",
			"chassis_no", "make_type", "model", "model_year", "device_name", "device_serial_no",
			"device_model", "device_manufacturer", "configured", "status", "device_details.created_on").
			From("device_details").
			LeftJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = device_details.vehicle_id")).
			Where(dbx.NewExp("device_details.vehicle_id > 0")).
			Where(dbx.HashExp{"vehicle_details.user_id": uid}).
			OrderBy("id DESC").Offset(int64(offset)).Limit(int64(limit))
	} else {
		q = rs.Tx().Select("id", "device_id", "device_details.vehicle_id", "vehicle_reg_no", "chassis_no", "make_type", "model", "model_year", "device_name", "device_serial_no", "device_model", "device_manufacturer", "configured", "status", "device_details.created_on").
			From("device_details").
			LeftJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = device_details.vehicle_id")).
			Where(dbx.NewExp("device_details.vehicle_id > 0")).
			OrderBy("id DESC").Offset(int64(offset)).Limit(int64(limit))
	}

	// populate data row by row
	rows, err := q.Rows()
	for rows.Next() {
		rows.ScanStruct(&d)
		deviceid := strconv.Itoa(int(d.DeviceID))
		d.Positions, _ = FetchCurrentPosition(deviceid, start, stop)
		devices = append(devices, d)
	}

	return devices, err
}

// FetchCurrentPosition ...
func FetchCurrentPosition(deviceid string, start, stop int64) ([]models.DeviceData, error) {
	var deviceData []models.DeviceData

	keysList, err := app.ZRevRange("data:"+deviceid, start, stop)
	if err != nil {
		fmt.Println("Getting Keys Failed : " + err.Error())
	}

	for i := 0; i < len(keysList); i++ {

		if keysList[i] != "0" {
			var deserializedValue models.DeviceData
			json.Unmarshal([]byte(keysList[i]), &deserializedValue)
			deviceData = append(deviceData, deserializedValue)
		}

	}

	return deviceData, err
}

// CountConfiguredDevices returns the number of the device configuration records in the database.
func (dao *DeviceDAO) CountConfiguredDevices(rs app.RequestScope, vehicleid int, deviceid int64) (int, error) {
	var count int
	query := "SELECT COUNT(device_id) FROM vehicle_configuration"
	if vehicleid > 0 && deviceid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "' AND device_id = '" + strconv.FormatInt(deviceid, 10) + "'"
	} else if vehicleid > 0 {
		query += " WHERE vehicle_id = '" + strconv.Itoa(vehicleid) + "'"
	} else if deviceid > 0 {
		query += " WHERE device_id = '" + strconv.FormatInt(deviceid, 10) + "'"
	} else {
		query += " WHERE status=1 "
	}
	q := rs.Tx().NewQuery(query)
	err := q.Row(&count)

	return count, err
}

// ConfiguredDevices retrieves the device records with the specified offset and limit from the database.
func (dao *DeviceDAO) ConfiguredDevices(rs app.RequestScope, offset, limit, vehicleid int, deviceid int64) ([]models.DeviceConfiguration, error) {
	devices := []models.DeviceConfiguration{}
	query := "SELECT conf_id, vc.device_id, vc.vehicle_id, vehicle_string_id AS device_name, JSON_VALUE(data, '$.sim_imei') AS sim_imei, vc.created_on, vc.status AS status, "
	query += " JSON_VALUE(data, '$.device_detail.chasis_no') AS chassis_no, JSON_VALUE(data, '$.device_detail.make_type') AS make_type, JSON_VALUE(data, '$.device_detail.device_type') AS device_type, "
	query += " vc.serial_no AS serial_no, JSON_VALUE(data, '$.device_detail.preset_speed') AS preset_speed, JSON_VALUE(data, '$.device_detail.set_frequency') AS set_frequency, "
	query += " vc.fitting_date AS fitting_date, DATE_ADD(DATE_ADD(COALESCE(renewal_date, vd.created_on), INTERVAL -1 DAY), INTERVAL 1 YEAR) AS expiry_date, JSON_VALUE(data, '$.device_detail.fitting_center') AS fitting_center, "
	query += " vc.certificate_no AS certificate, JSON_VALUE(data, '$.device_detail.email_address') AS email_address, JSON_VALUE(data, '$.device_detail.agent_phone') AS agent_phone, "
	query += " JSON_VALUE(data, '$.device_detail.agent_location') AS agent_location, JSON_VALUE(data, '$.device_detail.owner_name') AS owner_name, JSON_VALUE(data, '$.device_detail.owner_phone_number') AS owner_phone_number,"
	query += " COALESCE(dd.status, 0) AS device_status, COALESCE(dd.status_reason, 'Device record does not exit') AS reason FROM vehicle_configuration AS vc"
	query += " LEFT JOIN device_details AS dd ON (dd.device_id = vc.device_id)"

	if vehicleid > 0 && deviceid > 0 {
		query += " WHERE vc.vehicle_id = '" + strconv.Itoa(vehicleid) + "' AND device_id = '" + strconv.FormatInt(deviceid, 10) + "'"
	} else if vehicleid > 0 {
		query += " WHERE vc.vehicle_id = '" + strconv.Itoa(vehicleid) + "'"
	} else if deviceid > 0 {
		query += " WHERE vc.device_id = '" + strconv.FormatInt(deviceid, 10) + "'"
	} else {
		query += " WHERE vc.status=1 "
	}
	query += " ORDER BY conf_id DESC LIMIT " + strconv.Itoa(offset) + ", " + strconv.Itoa(limit)
	q := rs.Tx().NewQuery(query)
	err := q.All(&devices)

	return devices, err
}

// CountSearches ///
func (dao *DeviceDAO) CountSearches(rs app.RequestScope, searchterm string) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("device_details").Where(dbx.HashExp{"device_serial_no": searchterm}).Row(&count)
	return count, err
}

// SearchDevices retrieves the device records with the specified offset and limit from the database.
func (dao *DeviceDAO) SearchDevices(rs app.RequestScope, searchterm string, offset, limit int) ([]models.Devices, error) {
	devices := []models.Devices{}
	err := rs.Tx().Select("id", "device_id", "vehicle_id", "device_name", "device_details.company_id",
		"COALESCE(company_name, '') AS company_name", "certificate_no", "device_serial_no", "device_model", "device_manufacturer",
		"sim_serial_no", "sim_number", "motherboard_no", "technician", "status", "configured", "configuration_date", "status_reason", "created_on").
		LeftJoin("companies", dbx.NewExp("companies.company_id = device_details.company_id")).
		From("device_details").Where(dbx.HashExp{"device_details.device_serial_no": searchterm}).
		OrderBy("id DESC").Offset(int64(offset)).Limit(int64(limit)).All(&devices)

	return devices, err
}
