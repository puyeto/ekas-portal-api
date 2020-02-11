package daos

import (
	"encoding/json"
	"strconv"
	"time"

	// "time"
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// VehicleDAO persists vehicle data in database
type VehicleDAO struct{}

// NewVehicleDAO creates a new VehicleDAO
func NewVehicleDAO() *VehicleDAO {
	return &VehicleDAO{}
}

// GetVehicleByStrID ...
func (dao *VehicleDAO) GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error) {
	var vdetails models.VehicleConfigDetails
	query := "SELECT conf_id, vc.device_id, vd.user_id, vd.vehicle_id, vehicle_status, send_to_ntsa AS ntsa_show, vc.owner_id, fitter_id, notification_email, notification_no, data FROM vehicle_configuration AS vc "
	query += " LEFT JOIN vehicle_details AS vd ON (vd.vehicle_string_id = vc.vehicle_string_id) "
	query += " WHERE status=1 AND vc.vehicle_string_id='" + strid + "' LIMIT 1"
	q := rs.Tx().NewQuery(query)
	err := q.Row(&vdetails.ConfigID, &vdetails.DeviceID, &vdetails.UserID, &vdetails.VehicleID, &vdetails.VehicleStatus, &vdetails.NTSAShow, &vdetails.OwnerID, &vdetails.FitterID, &vdetails.NotificationEmail, &vdetails.NotificationNO, &vdetails.Data)
	return &vdetails, err
}

// GetConfigurationDetails ...
func (dao *VehicleDAO) GetConfigurationDetails(rs app.RequestScope, vehicleid, deviceid int) (*models.VehicleConfigDetails, error) {
	var vdetails models.VehicleConfigDetails
	query := "SELECT conf_id, vc.device_id, vc.vehicle_id, vc.owner_id, fitter_id, notification_email, notification_no, data FROM vehicle_configuration AS vc "
	query += " LEFT JOIN vehicle_details ON (vehicle_details.vehicle_id = vc.vehicle_id) "
	if deviceid > 0 && vehicleid > 0 {
		query += " WHERE status=1 AND vc.vehicle_id='" + strconv.Itoa(vehicleid) + "' AND vc.device_id='" + strconv.Itoa(deviceid) + "' "
	} else if deviceid > 0 {
		query += " WHERE status=1 AND vc.device_id='" + strconv.Itoa(deviceid) + "' "
	} else {
		query += " WHERE status=1 AND vc.vehicle_id='" + strconv.Itoa(vehicleid) + "' "
	}

	query += " LIMIT 1"
	q := rs.Tx().NewQuery(query)
	err := q.Row(&vdetails.ConfigID, &vdetails.DeviceID, &vdetails.VehicleID, &vdetails.OwnerID, &vdetails.FitterID, &vdetails.NotificationEmail, &vdetails.NotificationNO, &vdetails.Data)
	return &vdetails, err
}

// GetOverspeedByDeviceID ...
func (dao *VehicleDAO) GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.TripData, error) {
	tdetails := []models.TripData{}
	err := app.SecondDBCon.Select("trip_id", "device_id", "data_date", "speed").From("data_" + deviceid).
		OrderBy("trip_id DESC").Offset(int64(offset)).Limit(int64(limit)).
		Where(dbx.And(dbx.HashExp{"device_id": deviceid}, dbx.NewExp("speed>84"))).
		All(&tdetails)
	return tdetails, err
}

// CountOverspeed returns the number of overspeed records in the database.
func (dao *VehicleDAO) CountOverspeed(rs app.RequestScope, deviceid string) (int, error) {
	var cnt int
	// check if table exist
	err := app.SecondDBCon.NewQuery("SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'ekas_portal_data') AND (TABLE_NAME = 'data_" + deviceid + "')").Row(&cnt)
	if cnt == 0 {
		return cnt, err
	}

	var count int
	err = app.SecondDBCon.Select("COUNT(*)").From("data_" + deviceid).
		Where(dbx.And(dbx.HashExp{"device_id": deviceid}, dbx.NewExp("speed>84"))).
		Row(&count)
	return count, err
}

// GetViolationsByDeviceID ...
func (dao *VehicleDAO) GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.TripData, error) {
	tdetails := []models.TripData{}
	var query = "disconnect>0"
	if reason == "failsafe" {
		query = "failsafe>0"
	}

	err := app.SecondDBCon.Select("trip_id", "device_id", "data_date", "failsafe", "disconnect").From("data_" + deviceid).
		OrderBy("trip_id DESC").Offset(int64(offset)).Limit(int64(limit)).
		Where(dbx.And(dbx.HashExp{"device_id": deviceid}, dbx.NewExp(query))).
		All(&tdetails)
	return tdetails, err
}

// CountViolations returns the number of violation records in the database.
func (dao *VehicleDAO) CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error) {
	var count int
	var query = "disconnect>0"
	if reason == "failsafe" {
		query = "failsafe>0"
	}

	var cnt int
	// check if table exist
	err := app.SecondDBCon.NewQuery("SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'ekas_portal_data') AND (TABLE_NAME = 'data_" + deviceid + "')").Row(&cnt)
	if cnt == 0 {
		return cnt, err
	}

	err = app.SecondDBCon.Select("COUNT(*)").From("data_" + deviceid).
		Where(dbx.And(dbx.HashExp{"device_id": deviceid}, dbx.NewExp(query))).
		Row(&count)
	return count, err
}

// SearchVehicles ...
// qtype can be ntsa or ...
func (dao *VehicleDAO) SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int, qtype string) ([]models.SearchDetails, error) {
	tdetails := []models.SearchDetails{}

	q := rs.Tx().Select("DISTINCT(vehicle_configuration.vehicle_id) AS vehicle_name", "data").
		From("vehicle_configuration").LeftJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = vehicle_configuration.vehicle_id"))
	if qtype == "ntsa" {
		q.Where(dbx.Or(dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"vehicle_configuration.vehicle_string_id": searchterm}, dbx.NewExp("send_to_ntsa=1")),
			dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"device_id": searchterm}, dbx.NewExp("send_to_ntsa=1"))))
	} else {
		q.Where(dbx.Or(dbx.And(dbx.NewExp("status=1"), dbx.Like("vehicle_configuration.vehicle_string_id", searchterm)),
			dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"device_id": searchterm})))

	}
	err := q.OrderBy("vehicle_configuration.vehicle_id DESC").All(&tdetails)
	return tdetails, err
}

// CountSearches ///
func (dao *VehicleDAO) CountSearches(rs app.RequestScope, searchterm, qtype string) (int, error) {
	var count int
	q := rs.Tx().Select("COUNT(*)").From("vehicle_configuration").
		InnerJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = vehicle_configuration.vehicle_id"))
	if qtype == "ntsa" {
		q.Where(dbx.Or(dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"vehicle_configuration.vehicle_string_id": searchterm}, dbx.NewExp("send_to_ntsa=1")),
			dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"device_id": searchterm}, dbx.NewExp("send_to_ntsa=1"))))
	} else {
		q.Where(dbx.Or(dbx.And(dbx.NewExp("status=1"), dbx.Like("vehicle_configuration.vehicle_string_id", searchterm)),
			dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"device_id": searchterm})))

	}
	err := q.Row(&count)
	return count, err
}

// CountTripRecords returns the number of trip records in the database.
func (dao *VehicleDAO) CountTripRecords(rs app.RequestScope, deviceid string) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("trip_data").
		Where(dbx.HashExp{"device_id": deviceid}).
		Row(&count)
	return count, err
}

// ListRecentViolations ...
func (dao *VehicleDAO) ListRecentViolations(rs app.RequestScope, offset, limit int, uid string) ([]models.CurrentViolations, error) {
	tdetails := []models.CurrentViolations{}

	query := "SELECT cv.device_id, user_id, name, overspeed_trip_data, overspeed_speed, COALESCE(td.data_date, '') AS overspeed_date, "
	query += " disconnect_trip_data, disconnect_trip_speed, COALESCE(tdd.data_date, '') AS disconnect_trip_date, "
	query += " failsafe_trip_data, failsafe_trip_speed, COALESCE(tdf.data_date, '') AS failsafe_trip_date, offline_trip_data, "
	query += " offline_trip_speed, COALESCE(tdo.data_date, '') AS offline_trip_date FROM current_violations AS cv "
	query += " LEFT JOIN trip_data AS td ON (td.trip_id = overspeed_trip_data) "
	query += " LEFT JOIN trip_data AS tdd ON (tdd.trip_id = disconnect_trip_data) "
	query += " LEFT JOIN trip_data AS tdf ON (tdf.trip_id = failsafe_trip_data) "
	query += " LEFT JOIN trip_data AS tdo ON (tdo.trip_id = offline_trip_data) "
	if uid != "0" {
		query += " WHERE user_id=" + uid
	}
	query += " ORDER BY cv.created_on DESC LIMIT " + strconv.Itoa(limit)

	//OrderBy("cv.created_on DESC").Offset(int64(offset)).Limit(int64(limit))
	err := rs.Tx().NewQuery(query).All(&tdetails)
	if err != nil {
		return tdetails, err
	}

	return tdetails, nil
}

// ----------------------------Add / Update Vehicle------------------------------------

// CreateVehicle saves a new vehicle record in the database.
func (dao *VehicleDAO) CreateVehicle(rs app.RequestScope, v *models.VehicleDetails) (uint32, error) {
	exists, _ := dao.VehicleExistsByStringID(rs, v.VehicleStringID)
	if exists == 1 {
		vehid, err := dao.UpdateVehicle(rs, v)
		return vehid, err
	}

	err := rs.Tx().Model(v).Insert("VehicleID", "UserID", "VehicleStringID", "VehicleRegNo", "ChassisNo", "MakeType", "NotificationEmail", "NotificationNO")
	return v.VehicleID, err
}

// UpdateVehicle ....
func (dao *VehicleDAO) UpdateVehicle(rs app.RequestScope, v *models.VehicleDetails) (uint32, error) {
	var vehID uint32
	q := rs.Tx().NewQuery("SELECT vehicle_id FROM vehicle_details WHERE vehicle_string_id='" + v.VehicleStringID + "' LIMIT 1")
	err := q.Row(&vehID)

	_, err = rs.Tx().Update("vehicle_details", dbx.Params{
		"user_id":            v.UserID,
		"vehicle_string_id":  v.VehicleStringID,
		"vehicle_reg_no":     v.VehicleRegNo,
		"chassis_no":         v.ChassisNo,
		"make_type":          v.MakeType,
		"notification_email": v.NotificationEmail,
		"notification_no":    v.NotificationNO},
		dbx.HashExp{"vehicle_id": vehID}).Execute()
	return vehID, err
}

// VehicleExists check vehicle if exists by vehicle id...
func (dao *VehicleDAO) VehicleExists(rs app.RequestScope, id uint32) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_details WHERE vehicle_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// VehicleExistsByStringID Check vehicle if exists by string id...
func (dao *VehicleDAO) VehicleExistsByStringID(rs app.RequestScope, strID string) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_details WHERE vehicle_string_id='" + strID + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// ----------------------Add / update config data----------------------

// VehicleExistsConfigurationByStringID Check if vehicle exists under vehicle_configuration by vehicle string id...
func (dao *VehicleDAO) VehicleExistsConfigurationByStringID(rs app.RequestScope, strID string) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_details WHERE vehicle_string_id='" + strID + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

//CreateConfiguration Add configuartion details to db
func (dao *VehicleDAO) CreateConfiguration(rs app.RequestScope, cd *models.Vehicle, ownerid uint32, fitterid uint32, vehicleid uint32, vehstringid string) error {
	// Delete Previous Configuration
	_, err := rs.Tx().Delete("vehicle_configuration", dbx.HashExp{"vehicle_string_id": vehstringid}).Execute()
	if err != nil {
		return err
	}

	a, _ := json.Marshal(cd)
	_, err = rs.Tx().Insert("vehicle_configuration", dbx.Params{
		"user_id":           cd.UserID,
		"device_id":         cd.GovernorDetails.DeviceID,
		"vehicle_id":        vehicleid,
		"owner_id":          ownerid,
		"fitter_id":         fitterid,
		"vehicle_string_id": strings.ToLower(strings.Replace(cd.DeviceDetails.RegistrationNO, " ", "", -1)),
		"fitting_date":      cd.DeviceDetails.FittingDate,
		"frequency":         cd.DeviceDetails.SetFrequency,
		"speed":             cd.DeviceDetails.PresetSpeed,
		"speed_source":      cd.DeviceDetails.SpeedSource,
		"fail_safe":         cd.GovernorDetails.FailSafe,
		"apn":               cd.GovernorDetails.APN,
		"data":              string(a)}).Execute()
	return err
}

// UpdateConfigurationStatus ...
func (dao *VehicleDAO) UpdateConfigurationStatus(rs app.RequestScope, configid uint32, status int8) error {
	_, err := rs.Tx().Update("vehicle_configuration", dbx.Params{
		"status": status},
		dbx.HashExp{"conf_id": configid}).Execute()
	return err
}

// UpdatDeviceConfigurationStatus ...
func (dao *VehicleDAO) UpdatDeviceConfigurationStatus(rs app.RequestScope, deviceid int64, vehicleid uint32) error {
	t := time.Now()
	currentDate := t.Format("2006-01-02 15:04:05")

	_, err := rs.Tx().Update("device_details", dbx.Params{
		"configuration_date": currentDate, "configured": 1, "vehicle_id": vehicleid},
		dbx.HashExp{"device_id": deviceid}).Execute()
	return err
}

// CountTripDataByDeviceID returns the number of trip records in the database.
func (dao *VehicleDAO) CountTripDataByDeviceID(deviceid string) (int, error) {
	var count int
	var cnt int

	// check if table exist
	err := app.SecondDBCon.NewQuery("SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'ekas_portal_data') AND (TABLE_NAME = 'data_" + deviceid + "')").Row(&cnt)
	if cnt == 0 {
		return count, nil
	}

	err = app.SecondDBCon.Select("COUNT(*)").From("data_" + deviceid).
		Row(&count)
	return count, err
}

// GetTripDataByDeviceID ...
func (dao *VehicleDAO) GetTripDataByDeviceID(deviceid string, offset, limit int, orderby string) ([]models.DeviceData, error) {
	ddetails := []models.DeviceData{}
	var cnt int
	// check if table exist
	err := app.SecondDBCon.NewQuery("SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'ekas_portal_data') AND (TABLE_NAME = 'data_" + deviceid + "')").Row(&cnt)
	if cnt == 0 {
		return ddetails, nil
	}

	err = app.SecondDBCon.Select("device_id", "data_date AS date_time", "speed AS ground_speed", "latitude", "longitude", "date_time_stamp").From("data_" + deviceid).
		OrderBy("date_time_stamp " + orderby).Offset(int64(offset)).Limit(int64(limit)).All(&ddetails)
	return ddetails, err
}

// CountTripRecordsBtwDates returns the number of trip records between dates in the database.
func (dao *VehicleDAO) CountTripRecordsBtwDates(deviceid string, from, to int64) (int, error) {
	var count int
	var cnt int
	// check if table exist
	err := app.SecondDBCon.NewQuery("SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'ekas_portal_data') AND (TABLE_NAME = 'data_" + deviceid + "')").Row(&cnt)
	if cnt == 0 {
		return count, nil
	}

	err = app.SecondDBCon.Select("COUNT(*)").From("data_" + deviceid).
		Where(dbx.And(dbx.Between("date_time_stamp", from, to), dbx.HashExp{"device_id": deviceid})).
		Row(&count)
	return count, err
}

// GetTripDataByDeviceIDBtwDates ...
func (dao *VehicleDAO) GetTripDataByDeviceIDBtwDates(deviceid string, offset, limit int, from, to int64) ([]models.DeviceData, error) {
	tdetails := []models.DeviceData{}
	var cnt int

	// check if table exist
	err := app.SecondDBCon.NewQuery("SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'ekas_portal_data') AND (TABLE_NAME = 'data_" + deviceid + "')").Row(&cnt)
	if cnt == 0 {
		return tdetails, nil
	}

	err = app.SecondDBCon.Select("device_id", "data_date AS date_time", "speed AS ground_speed", "latitude", "longitude").From("data_" + deviceid).
		Where(dbx.And(dbx.Between("date_time_stamp", from, to), dbx.HashExp{"device_id": deviceid})).
		OrderBy("date_time_stamp DESC").Offset(int64(offset)).Limit(int64(limit)).All(&tdetails)
	return tdetails, err
}

// CreateDevice saves a new device record in the database.
// The Device.Id field will be populated with an automatically generated ID upon successful saving.
func (dao *VehicleDAO) CreateDevice(rs app.RequestScope, device *models.Devices) error {
	var exists int
	strID := strconv.FormatInt(device.DeviceID, 10)
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM device_details WHERE device_id='" + strID + "' LIMIT 1) AS exist")
	err := q.Row(&exists)

	if exists == 1 {
		_, err = rs.Tx().Update("device_details", dbx.Params{
			"device_serial_no": strings.ToUpper(device.DeviceSerialNo),
			"sim_serial_no":    device.SimSerialNo,
			"sim_number":       device.SimNumber,
			"motherboard_no":   device.MotherboardNO,
			"technician":       device.Technician,
			"configured":       device.Configured,
			"note":             device.Note,
		}, dbx.HashExp{"device_id": device.DeviceID}).Execute()
		return err
	}

	_, err = rs.Tx().Insert("device_details", dbx.Params{
		"device_id":           device.DeviceID,
		"device_name":         device.DeviceName,
		"device_serial_no":    strings.ToUpper(device.DeviceSerialNo),
		"device_model":        strings.ToUpper(device.DeviceModelNo),
		"device_manufacturer": strings.ToUpper(device.DeviceManufacturer),
		"sim_serial_no":       device.SimSerialNo,
		"sim_number":          device.SimNumber,
		"motherboard_no":      device.MotherboardNO,
		"technician":          device.Technician,
		"configured":          device.Configured,
		"note":                device.Note,
	}).Execute()

	return err
}
