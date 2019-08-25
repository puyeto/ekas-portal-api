package daos

import (
	"encoding/json"
	"fmt"
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
	query := "SELECT conf_id, vc.device_id, vc.vehicle_id, vc.owner_id, fitter_id, notification_email, notification_no, data FROM vehicle_configuration AS vc "
	query += " LEFT JOIN vehicle_details ON (vehicle_details.vehicle_id = vc.vehicle_id) "
	query += " WHERE status=1 AND vc.vehicle_string_id='" + strid + "' LIMIT 1"
	q := rs.Tx().NewQuery(query)
	err := q.Row(&vdetails.ConfigID, &vdetails.DeviceID, &vdetails.VehicleID, &vdetails.OwnerID, &vdetails.FitterID, &vdetails.NotificationEmail, &vdetails.NotificationNO, &vdetails.Data)
	fmt.Println(vdetails)
	return &vdetails, err
}

// GetOverspeedByDeviceID ...
func (dao *VehicleDAO) GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.TripData, error) {
	tdetails := []models.TripData{}
	err := rs.Tx().Select("trip_id", "device_id", "data_date", "speed").
		OrderBy("trip_id DESC").Offset(int64(offset)).Limit(int64(limit)).
		// Where(dbx.HashExp{"device_id": deviceid, "speed>": 80}).
		Where(dbx.And(dbx.HashExp{"device_id": deviceid}, dbx.NewExp("speed>80"))).
		All(&tdetails)
	return tdetails, err
}

// GetViolationsByDeviceID ...
func (dao *VehicleDAO) GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.TripData, error) {
	tdetails := []models.TripData{}
	var query = "disconnect>0"
	if reason == "failsafe" {
		query = "failsafe>0"
	}
	err := rs.Tx().Select("trip_id", "device_id", "data_date", "failsafe", "disconnect").
		OrderBy("trip_id DESC").Offset(int64(offset)).Limit(int64(limit)).
		Where(dbx.And(dbx.HashExp{"device_id": deviceid}, dbx.NewExp(query))).
		All(&tdetails)
	return tdetails, err
}

// SearchVehicles ...
func (dao *VehicleDAO) SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int) ([]models.SearchDetails, error) {
	tdetails := []models.SearchDetails{}

	err := rs.Tx().Select("DISTINCT(vehicle_id) AS vehicle_name", "data").
		From("vehicle_configuration").Offset(int64(offset)).Limit(int64(limit)).
		// InnerJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = vehicle_configuration.vehicle_id")).
		Where(dbx.Like("vehicle_string_id", searchterm)).
		All(&tdetails)
	return tdetails, err
}

// CountTripRecords returns the number of trip records in the database.
func (dao *VehicleDAO) CountTripRecords(rs app.RequestScope, deviceid string) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("trip_data").
		Where(dbx.HashExp{"device_id": deviceid}).
		Row(&count)
	return count, err
}

// CountSearches ///
func (dao *VehicleDAO) CountSearches(rs app.RequestScope, searchterm string) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("vehicle_configuration").
		InnerJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = vehicle_configuration.vehicle_id")).
		Where(dbx.Like("vehicle_configuration.vehicle_string_id", searchterm)).
		Row(&count)
	return count, err
}

// CountOverspeed returns the number of overspeed records in the database.
func (dao *VehicleDAO) CountOverspeed(rs app.RequestScope, deviceid string) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("trip_data").
		Where(dbx.And(dbx.HashExp{"device_id": deviceid}, dbx.NewExp("speed>80"))).
		Row(&count)
	return count, err
}

// CountViolations returns the number of violation records in the database.
func (dao *VehicleDAO) CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error) {
	var count int
	var query = "disconnect>0"
	if reason == "failsafe" {
		query = "failsafe>0"
	}
	err := rs.Tx().Select("COUNT(*)").From("trip_data").
		Where(dbx.And(dbx.HashExp{"device_id": deviceid}, dbx.NewExp(query))).
		Row(&count)
	return count, err
}

// CountTripRecordsBtwDates returns the number of trip records between dates in the database.
func (dao *VehicleDAO) CountTripRecordsBtwDates(rs app.RequestScope, deviceid string, from string, to string) (int, error) {
	// formatedFrom := from.Format("2006-01-02 15:04:05")
	// formatedTo := to.Format("2006-01-02 15:04:05")
	var count int
	err := rs.Tx().Select("COUNT(*)").From("trip_data").
		Where(dbx.And(dbx.Between("data_date", from, to), dbx.HashExp{"device_id": deviceid})).
		Row(&count)
	return count, err
}

// FetchAllTripsBetweenDates ...
func (dao *VehicleDAO) FetchAllTripsBetweenDates(rs app.RequestScope, deviceid string, offset, limit int, from string, to string) ([]models.TripData, error) {
	tdetails := []models.TripData{}
	err := rs.Tx().Select("trip_id", "device_id", "data_date", "speed", "longitude", "latitude").
		OrderBy("trip_id DESC").Offset(int64(offset)).Limit(int64(limit)).
		Where(dbx.And(dbx.Between("data_date", from, to), dbx.HashExp{"device_id": deviceid})).All(&tdetails)
	return tdetails, err
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
func (dao *VehicleDAO) CreateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	exists, _ := dao.VehicleExists(rs, v.VehicleID)
	if exists == 1 {
		return dao.UpdateVehicle(rs, v)
	}

	return rs.Tx().Model(v).Insert("VehicleID", "UserID", "VehicleStringID", "VehicleRegNo", "ChassisNo", "MakeType", "NotificationEmail", "NotificationNO")
}

// UpdateVehicle ....
func (dao *VehicleDAO) UpdateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	_, err := rs.Tx().Update("vehicle_details", dbx.Params{
		"user_id":            v.UserID,
		"vehicle_string_id":  v.VehicleStringID,
		"vehicle_reg_no":     v.VehicleRegNo,
		"chassis_no":         v.ChassisNo,
		"make_type":          v.MakeType,
		"notification_email": v.NotificationEmail,
		"notification_no":    v.NotificationNO},
		dbx.HashExp{"vehicle_id": v.VehicleID}).Execute()
	return err
}

// VehicleExists ...
func (dao *VehicleDAO) VehicleExists(rs app.RequestScope, id uint32) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_details WHERE vehicle_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// ----------------------Add / update config data----------------------

//CreateConfiguration Add configuartion details to db
func (dao *VehicleDAO) CreateConfiguration(rs app.RequestScope, cd *models.Vehicle, ownerid uint32, fitterid uint32, vehicleid uint32) error {
	a, _ := json.Marshal(cd)
	_, err := rs.Tx().Insert("vehicle_configuration", dbx.Params{
		"user_id":           cd.UserID,
		"device_id":         cd.GovernorDetails.DeviceID,
		"vehicle_id":        vehicleid,
		"owner_id":          ownerid,
		"fitter_id":         fitterid,
		"vehicle_string_id": strings.ToLower(strings.Replace(cd.DeviceDetails.RegistrationNO, " ", "", -1)),
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
func (dao *VehicleDAO) UpdatDeviceConfigurationStatus(rs app.RequestScope, deviceid uint32, vehicleid uint32) error {
	t := time.Now()
	currentDate := t.Format("2006-01-02 15:04:05")

	_, err := rs.Tx().Update("device_details", dbx.Params{
		"configuration_date": currentDate, "configured": 1, "vehicle_id": vehicleid},
		dbx.HashExp{"device_id": deviceid}).Execute()
	return err
}
