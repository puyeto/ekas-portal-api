package daos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	// "time"
	"strings"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VehicleDAO persists vehicle data in database
type VehicleDAO struct{}

// NewVehicleDAO creates a new VehicleDAO
func NewVehicleDAO() *VehicleDAO {
	return &VehicleDAO{}
}

// GetVehicleName ...
func (dao *VehicleDAO) GetVehicleName(rs app.RequestScope, deviceid int) models.VDetails {
	var vd models.VDetails
	query := "SELECT send_to_ntsa, device_id, sim_no, vehicle_reg_no, json_value(data, '$.device_detail.owner_name'), json_value(data, '$.device_detail.owner_phone_number') "
	query += " FROM vehicle_configuration "
	query += " LEFT JOIN vehicle_details AS vd ON (vd.vehicle_id = vehicle_configuration.vehicle_id) "
	query += " WHERE device_id='" + strconv.Itoa(deviceid) + "' LIMIT 1"
	rs.Tx().NewQuery(query).Row(&vd.SendToNTSA, &vd.DeviceID, &vd.DeviceSIMNo, &vd.Name, &vd.VehicleOwner, &vd.OwnerTel)

	return vd
}

// GetVehicleByStrID ...
func (dao *VehicleDAO) GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error) {
	var vdetails models.VehicleConfigDetails

	vehicleid, err := GetVehicleIDByStrID(rs, strid)
	if err != nil {
		return &vdetails, err
	}

	query := "SELECT conf_id, vc.device_id, vd.user_id, COALESCE(CONCAT(u.first_name , ' ' , u.last_name), '') AS fitter, vd.vehicle_id, vd.vehicle_reg_no, vehicle_status, send_to_ntsa AS ntsa_show, vc.owner_id, "
	query += " fitter_id, notification_email, notification_no, COALESCE(JSON_VALUE(data, '$.device_detail.sim_no'), '') AS sim_no, serial_no, last_seen, COALESCE(vr.renewal_date, vd.created_on) AS renewal_date, renew, "
	query += " vd.created_on, DATE_ADD(DATE_ADD(COALESCE(vr.renewal_date, vd.created_on), INTERVAL -1 DAY), INTERVAL 1 YEAR) AS expiry_date, device_status, vd.sacco_id, data FROM vehicle_configuration AS vc "
	query += " LEFT JOIN vehicle_details AS vd ON (vd.vehicle_id = vc.vehicle_id) "
	query += " LEFT JOIN auth_users AS u ON (u.auth_user_id = vd.user_id) "
	query += " LEFT JOIN (SELECT * FROM vehicle_renewals WHERE vehicle_id='" + strconv.Itoa(vehicleid) + "' ORDER BY id DESC LIMIT 1) AS vr ON (vr.vehicle_id = vd.vehicle_id) "
	query += " WHERE vc.status=1 AND vc.vehicle_id='" + strconv.Itoa(vehicleid) + "' LIMIT 1"
	q := rs.Tx().NewQuery(query)
	err = q.Row(&vdetails.ConfigID, &vdetails.DeviceID, &vdetails.UserID, &vdetails.Fitter, &vdetails.VehicleID, &vdetails.VehicleRegistration, &vdetails.VehicleStatus, &vdetails.NTSAShow, &vdetails.OwnerID, &vdetails.FitterID, &vdetails.NotificationEmail, &vdetails.NotificationNO, &vdetails.SimNO, &vdetails.SerialNo, &vdetails.LastSeen, &vdetails.RenewalDate, &vdetails.Renew, &vdetails.CreatedOn, &vdetails.ExpiryDate, &vdetails.DeviceStatus, &vdetails.SaccoID, &vdetails.Data)
	return &vdetails, err
}

func GetVehicleIDByStrID(rs app.RequestScope, strid string) (int, error) {
	var vid int
	err := rs.Tx().Select("vehicle_id").Where(dbx.HashExp{"vehicle_string_id": strid}).From("vehicle_details").OrderBy("vehicle_id DESC").Limit(1).Row(&vid)
	return vid, err
}

// GetConfigurationDetails ...
func (dao *VehicleDAO) GetConfigurationDetails(rs app.RequestScope, vehicleid int, deviceid int64) (*models.VehicleConfigDetails, error) {
	var vdetails models.VehicleConfigDetails
	var vid = vehicleid
	if vid == 0 {
		// get vehicle details
		rs.Tx().Select("vehicle_id").Where(dbx.HashExp{"device_id": deviceid}).From("vehicle_configuration").OrderBy("conf_id DESC").Limit(1).Row(&vid)
	}

	query := "SELECT conf_id, vc.device_id, vd.user_id, COALESCE(CONCAT(u.first_name , ' ' , u.last_name), '') AS fitter, vd.vehicle_id, vd.vehicle_reg_no, vehicle_status, send_to_ntsa AS ntsa_show, vc.owner_id, "
	query += " fitter_id, notification_email, notification_no, COALESCE(JSON_VALUE(data, '$.device_detail.sim_no'), '') AS sim_no, serial_no, last_seen, COALESCE(vr.renewal_date, vd.created_on) AS renewal_date, renew, "
	query += " vd.created_on, DATE_ADD(DATE_ADD(COALESCE(vr.renewal_date, vd.created_on), INTERVAL -1 DAY), INTERVAL 1 YEAR) AS expiry_date, device_status, data FROM vehicle_configuration AS vc "
	query += " LEFT JOIN vehicle_details AS vd ON (vd.vehicle_id = vc.vehicle_id) "
	query += " LEFT JOIN auth_users AS u ON (u.auth_user_id = vd.user_id) "
	query += " LEFT JOIN (SELECT * FROM vehicle_renewals WHERE vehicle_id='" + strconv.Itoa(vid) + "' ORDER BY id DESC LIMIT 1) AS vr ON (vr.vehicle_id = vd.vehicle_id) "

	if deviceid > 0 && vehicleid > 0 {
		query += " WHERE vc.status=1 AND vc.vehicle_id='" + strconv.Itoa(vehicleid) + "' AND vc.device_id='" + strconv.FormatInt(deviceid, 10) + "' "
	} else if deviceid > 0 {
		query += " WHERE vc.status=1 AND vc.device_id='" + strconv.FormatInt(deviceid, 10) + "' "
	} else {
		query += " WHERE vc.status=1 AND vc.vehicle_id='" + strconv.Itoa(vehicleid) + "' "
	}

	query += " LIMIT 1"
	q := rs.Tx().NewQuery(query)
	err := q.Row(&vdetails.ConfigID, &vdetails.DeviceID, &vdetails.UserID, &vdetails.Fitter, &vdetails.VehicleID, &vdetails.VehicleRegistration, &vdetails.VehicleStatus, &vdetails.NTSAShow, &vdetails.OwnerID, &vdetails.FitterID, &vdetails.NotificationEmail, &vdetails.NotificationNO, &vdetails.SimNO, &vdetails.SerialNo, &vdetails.LastSeen, &vdetails.RenewalDate, &vdetails.Renew, &vdetails.CreatedOn, &vdetails.ExpiryDate, &vdetails.DeviceStatus, &vdetails.Data)
	return &vdetails, err
}

// GetOverspeedByDeviceID ...
func (dao *VehicleDAO) GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.DeviceData, error) {
	findOptions := options.Find()
	findOptions.SetSort(map[string]int{"datetimestamp": -1})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	filter := bson.D{{Key: "groundspeed", Value: bson.D{{Key: "$gt", Value: 84}}}}
	return app.GetDeviceDataLogsMongo(deviceid, filter, findOptions)
}

// DeleteOverspeedsByDeviceID deletes an Record with the specified ID from the database (mongodb).
func (dao *VehicleDAO) DeleteOverspeedsByDeviceID(rs app.RequestScope, id uint64) (int, error) {
	filter := bson.D{{Key: "groundspeed", Value: bson.D{{Key: "$gt", Value: 81}}}}

	// Get collection
	collection := app.MongoDB.Collection("data_" + strconv.FormatInt(int64(id), 10))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(res.DeletedCount), nil
}

// DeleteFutureDataByDeviceID delete Records with with future from the database (mongodb).
func (dao *VehicleDAO) DeleteFutureDataByDeviceID(rs app.RequestScope, id uint64) (int, error) {
	now := time.Now()
	nowTimestamp := now.Unix()
	filter := bson.D{{Key: "datetimestamp", Value: bson.D{{Key: "$gt", Value: nowTimestamp}}}}

	// Get collection
	collection := app.MongoDB.Collection("data_" + strconv.FormatInt(int64(id), 10))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(res.DeletedCount), nil
}

// CountOverspeed returns the number of overspeed records in the database.
func (dao *VehicleDAO) CountOverspeed(rs app.RequestScope, deviceid string) (int, error) {
	filter := bson.D{{Key: "groundspeed", Value: bson.D{{Key: "$gt", Value: 84}}}}
	count, err := Count(deviceid, filter, nil)
	return int(count), err
}

// GetViolationsByDeviceID ...
func (dao *VehicleDAO) GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.DeviceData, error) {
	findOptions := options.Find()
	findOptions.SetSort(map[string]int{"datetimestamp": -1})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	filter := bson.D{}
	if reason == "failsafe" {
		filter = bson.D{{Key: "failsafe", Value: true}}
	} else {
		filter = bson.D{{Key: "disconnect", Value: true}}
	}
	return app.GetDeviceDataLogsMongo(deviceid, filter, findOptions)
}

// CountViolations returns the number of violation records in the database.
func (dao *VehicleDAO) CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error) {
	filter := bson.M{}
	if reason == "failsafe" {
		filter = bson.M{"failsafe": true}
	} else {
		filter = bson.M{"disconnect": true}
	}

	count, err := app.CountRecordsMongo("data_"+deviceid, filter, nil)
	return count, err

}

// SearchVehicles ...
// qtype can be ntsa or ...
func (dao *VehicleDAO) SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int, qtype string) ([]models.SearchDetails, error) {
	tdetails := []models.SearchDetails{}

	q := rs.Tx().Select("DISTINCT(vehicle_configuration.vehicle_id) AS vehicle_name", "vehicle_details.vehicle_id", "vehicle_configuration.device_id", "vehicle_details.vehicle_reg_no", "data").
		From("vehicle_configuration").LeftJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = vehicle_configuration.vehicle_id"))
	if qtype == "ntsa" {
		q.Where(dbx.Or(dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"vehicle_configuration.vehicle_string_id": searchterm}, dbx.NewExp("send_to_ntsa=1")),
			dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"vehicle_details.vehicle_status": 1}, dbx.HashExp{"device_id": searchterm}, dbx.NewExp("send_to_ntsa=1"))))
	} else {
		q.Where(dbx.Or(dbx.And(dbx.NewExp("status=1"), dbx.Like("vehicle_configuration.vehicle_string_id", searchterm)),
			dbx.And(dbx.NewExp("status=1"), dbx.Like("vehicle_configuration.serial_no", searchterm)),
			dbx.And(dbx.NewExp("status=1"), dbx.HashExp{"vehicle_details.vehicle_status": 1}, dbx.HashExp{"device_id": searchterm})))
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
			dbx.And(dbx.NewExp("status=1"), dbx.Like("vehicle_configuration.serial_no", searchterm)),
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
func (dao *VehicleDAO) CreateConfiguration(rs app.RequestScope, cd *models.Vehicle, ownerid uint32, fitterid uint32, vehicleid uint32) error {
	var vehiclestringid = strings.ToLower(strings.Replace(cd.DeviceDetails.RegistrationNO, " ", "", -1))

	// Delete Previous Configuration
	_, err := rs.Tx().Delete("vehicle_configuration", dbx.HashExp{"vehicle_string_id": vehiclestringid}).Execute()
	if err != nil {
		return err
	}

	rs.Tx().Insert("vehicle_devices", dbx.Params{
		"device_id":         cd.GovernorDetails.DeviceID,
		"vehicle_string_id": vehiclestringid,
		"vehicle_id":        vehicleid,
		"created_on":        time.Now()}).Execute()

	// rs.Tx().Update("vehicle_details", dbx.Params{
	// 	"send_to_ntsa": 0},
	// 	dbx.HashExp{"vehicle_id": vehicleid}).Execute()

	a, _ := json.Marshal(cd)
	_, err = rs.Tx().Insert("vehicle_configuration", dbx.Params{
		"user_id":           cd.UserID,
		"device_id":         cd.GovernorDetails.DeviceID,
		"vehicle_id":        vehicleid,
		"owner_id":          ownerid,
		"fitter_id":         fitterid,
		"vehicle_string_id": vehiclestringid,
		"fitting_date":      cd.DeviceDetails.FittingDate,
		"frequency":         cd.DeviceDetails.SetFrequency,
		"speed":             cd.DeviceDetails.PresetSpeed,
		"speed_source":      cd.DeviceDetails.SpeedSource,
		"fail_safe":         cd.GovernorDetails.FailSafe,
		"apn":               cd.GovernorDetails.APN,
		"serial_no":         cd.DeviceDetails.SerialNO,
		"certificate_no":    cd.DeviceDetails.Certificate,
		"sim_no":            cd.SimNO,
		"data":              string(a)}).Execute()
	return err
}

// CheckIfSerialNoExists ...
func (dao *VehicleDAO) CheckIfSerialNoExists(rs app.RequestScope, cd *models.Vehicle) error {
	var vehiclestringid = strings.ToLower(strings.Replace(cd.DeviceDetails.RegistrationNO, " ", "", -1))
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_configuration WHERE vehicle_string_id!='" + vehiclestringid + "' AND serial_no='" + cd.DeviceDetails.SerialNO + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	if err != nil {
		return errors.New("serial check failed try again")
	}
	if exists > 0 {
		return errors.New("serial No already exist with a different vehicle")
	}

	return nil
}

// CheckIfDeviceIDExists ...
func (dao *VehicleDAO) CheckIfDeviceIDExists(rs app.RequestScope, cd *models.Vehicle) error {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_configuration WHERE device_id='" + cd.GovernorDetails.DeviceID + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	if err != nil {
		return errors.New("Device ID check failed try again")
	}
	if exists > 0 {
		return errors.New("Device ID already exist. Please unplug your device and try again")
	}

	return nil
}

// CheckIfVehicleIsExpired ...
func (dao *VehicleDAO) CheckIfVehicleIsExpired(rs app.RequestScope, cd *models.Vehicle, daystoexpiry int) error {
	now := time.Now()

	// get date of expiry
	var vehiclestringid = strings.ToLower(strings.Replace(cd.DeviceDetails.RegistrationNO, " ", "", -1))
	var expiry time.Time
	// query := "SELECT DATE_ADD(DATE_ADD(COALESCE(renewal_date, created_on), INTERVAL -1 DAY), INTERVAL 1 YEAR) AS expiry_date FROM vehicle_details WHERE vehicle_string_id='" + vehiclestringid + "'"
	query := "SELECT DATE_ADD(DATE_ADD(COALESCE(vr.renewal_date, vd.created_on), INTERVAL -1 DAY), INTERVAL 1 YEAR) AS expiry_date "
	query += " FROM vehicle_details AS vd LEFT JOIN (SELECT * FROM vehicle_renewals WHERE vehicle_string_id='" + vehiclestringid + "' ORDER BY id DESC LIMIT 1) AS vr ON (vr.vehicle_id = vd.vehicle_id) WHERE vd.vehicle_string_id='" + vehiclestringid + "'"
	err := rs.Tx().NewQuery(query).Row(&expiry)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		return errors.New("expiry check failed try again")
	}

	diff := expiry.Sub(now).Hours() / 24
	if int(diff) <= daystoexpiry {
		return errors.New("This vehicle has expired or yet to expire")
	}

	return nil
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
	id, _ := strconv.ParseInt(deviceid, 10, 64)
	filter := bson.D{{Key: "deviceid", Value: id}}
	return Count(deviceid, filter, nil)
}

// Count returns the number of trip records in the database.
func Count(deviceid string, filter primitive.D, opts *options.FindOptions) (int, error) {
	app.CreateIndexMongo("data_" + deviceid)
	collection := app.MongoDB.Collection("data_" + deviceid)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // cancel when we are finished consuming integers
	count, err := collection.CountDocuments(ctx, filter, nil)
	return int(count), err
}

// GetTripDataByDeviceID ...
func (dao *VehicleDAO) GetTripDataByDeviceID(deviceid string, offset, limit int, orderby string) ([]models.DeviceData, error) {
	app.CreateIndexMongo(deviceid)
	findOptions := options.Find()
	if orderby == "desc" {
		// Sort by `price` field descending
		findOptions.SetSort(map[string]int{"datetimestamp": -1})
	}
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	filter := bson.D{}
	return app.GetDeviceDataLogsMongo(deviceid, filter, findOptions)
}

// CountTripRecordsBtwDates returns the number of trip records between dates in the database.
func (dao *VehicleDAO) CountTripRecordsBtwDates(deviceid string, from, to int64) (int, error) {
	// id, _ := strconv.Atoi(deviceid)
	// filter := bson.M{"deviceid": id}
	filter := bson.D{
		{Key: "datetimestamp", Value: bson.D{{Key: "$gte", Value: from}}},
		{Key: "datetimestamp", Value: bson.D{{Key: "$lte", Value: to}}},
	}
	count, err := Count(deviceid, filter, nil)
	fmt.Printf("count %v with error %v", count, err)
	return count, err
}

// GetTripDataByDeviceIDBtwDates ...
func (dao *VehicleDAO) GetTripDataByDeviceIDBtwDates(deviceid string, offset, limit int, from, to int64) ([]models.DeviceData, error) {
	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(map[string]int{"datetimestamp": -1})

	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	filter := bson.D{
		{Key: "datetimestamp", Value: bson.D{{Key: "$gte", Value: from}}},
		{Key: "datetimestamp", Value: bson.D{{Key: "$lte", Value: to}}},
	}
	return app.GetDeviceDataLogsMongo(deviceid, filter, findOptions)
}

// CountAllViolations returns the number of violation records in the database.
func (dao *VehicleDAO) CountAllViolations() (int, error) {
	count, err := app.CountRecordsMongo("current_violations", nil, nil)
	fmt.Printf("count %v with error %v", count, err)
	return count, err
}

// ListAllViolations ...
func (dao *VehicleDAO) ListAllViolations(rs app.RequestScope, offset, limit int) ([]models.CurrentViolations, error) {
	var vdetails []models.CurrentViolations
	findOptions := options.Find()
	findOptions.SetSort(map[string]int{"datetimeunix": -1})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	filter := bson.D{}

	cur, err := app.FindDataMongoDB("current_violations", filter, findOptions)
	if err != nil {
		return vdetails, err
	}
	for cur.Next(context.Background()) {
		item := models.CurrentViolations{}
		err := cur.Decode(&item)
		if err != nil {
			continue
		}
		vd := dao.GetVehicleName(rs, int(item.DeviceID))
		item.VehicleRegistration = vd.Name
		item.Data.Name = vd.Name
		item.VehicleOwner = vd.VehicleOwner
		item.OwnerTel = vd.OwnerTel
		if item.VehicleRegistration != "" {
			vdetails = append(vdetails, item)
		}
	}

	if err := cur.Err(); err != nil {
		return vdetails, err
	}

	return vdetails, err
}

// XMLListAllViolations ...
func (dao *VehicleDAO) XMLListAllViolations(rs app.RequestScope, offset, limit int) ([]models.XMLResults, error) {
	var vdetails []models.XMLResults
	findOptions := options.Find()
	findOptions.SetSort(map[string]int{"datetimeunix": -1})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	filter := bson.D{}

	cur, err := app.FindDataMongoDB("current_violations", filter, findOptions)
	if err != nil {
		return vdetails, err
	}
	for cur.Next(context.Background()) {
		var dData models.XMLResults
		item := models.CurrentViolations{}
		err := cur.Decode(&item)
		if err != nil {
			continue
		}
		vd := dao.GetVehicleName(rs, int(item.DeviceID))
		if vd.SendToNTSA == 0 {
			continue
		}
		dData.SerialNo = item.DeviceID
		dData.DateOfViolation = item.DateTime.Local().Format("2006-01-02 15:04:05")
		dData.VehicleRegistration = vd.Name
		dData.VehicleOwner = vd.VehicleOwner
		dData.OwnerTel = vd.OwnerTel
		dData.DeviceSIMNO = vd.DeviceSIMNo

		if item.Data.Failsafe {
			dData.ViolationType = "Signal Disconnect"
		} else if item.Data.Disconnect {
			dData.ViolationType = "Power Disconnect"
		} else if item.Data.Offline {
			dData.ViolationType = "Offline"
		} else {
			dData.ViolationType = "Overspeeding"
		}

		if dData.VehicleRegistration != "" {
			vdetails = append(vdetails, dData)
		}
	}

	return vdetails, cur.Err()
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

// GetFitterIDByAgentIDNo ...
func (dao *VehicleDAO) GetFitterIDByAgentIDNo(rs app.RequestScope, agentid int) uint32 {
	var aid uint32
	query := "SELECT company_id FROM companies WHERE contact_id='" + strconv.Itoa(agentid) + "' LIMIT 1"
	rs.Tx().NewQuery(query).Row(&aid)

	return aid
}
