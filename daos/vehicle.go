package daos

import (
	"encoding/json"
	"fmt"
	"strconv"

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
	query := "SELECT conf_id, vehicle_id, owner_id, fitter_id, data FROM vehicle_configuration "
	query += " WHERE status=1 AND vehicle_string_id='" + strid + "' LIMIT 1"
	q := rs.Tx().NewQuery(query)
	err := q.Row(&vdetails.ConfigID, &vdetails.VehicleID, &vdetails.OwnerID, &vdetails.FitterID, &vdetails.Data)
	fmt.Println(vdetails)
	return &vdetails, err
}

// GetTripDataByDeviceID ...
func (dao *VehicleDAO) GetTripDataByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.TripData, error) {
	tdetails := []models.TripData{}
	err := rs.Tx().Select("trip_id", "device_id", "data_date", "speed", "longitude", "latitude").
		OrderBy("trip_id DESC").Offset(int64(offset)).Limit(int64(limit)).
		Where(dbx.HashExp{"device_id": deviceid}).All(&tdetails)
	return tdetails, err
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

// CountTripRecords returns the number of trip records in the database.
func (dao *VehicleDAO) CountTripRecords(rs app.RequestScope, deviceid string) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("trip_data").
		Where(dbx.HashExp{"device_id": deviceid}).
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

// ----------------------------Add / Update Vehicle------------------------------------

// CreateVehicle saves a new vehicle record in the database.
func (dao *VehicleDAO) CreateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	if exists, _ := dao.VehicleExists(rs, v.VehicleID); exists == 1 {
		return dao.UpdateVehicle(rs, v)
	}

	return rs.Tx().Model(v).Insert("VehicleID", "VehicleStringID", "VehicleRegNo", "ChassisNo", "MakeType")
}

// UpdateVehicle ....
func (dao *VehicleDAO) UpdateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	_, err := rs.Tx().Update("vehicle_details", dbx.Params{
		"vehicle_string_id": v.VehicleStringID,
		"vehicle_reg_no":    v.VehicleRegNo,
		"chassis_no":        v.ChassisNo,
		"make_type":         v.MakeType},
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

// ------------------------Add / Update Owner-----------------------------------

// CreateVehicleOwner Add vehicle owner
func (dao *VehicleDAO) CreateVehicleOwner(rs app.RequestScope, vo *models.VehicleOwner) error {
	if exists, _ := dao.VehicleOwnerExists(rs, vo.OwnerID); exists == 1 {
		return dao.UpdateVehicleOwners(rs, vo)
	}
	return rs.Tx().Model(vo).Insert("OwnerID", "OwnerIDNo", "OwnerName", "OwnerEmail", "OwnerPhone")
}

// UpdateVehicleOwners ....
func (dao *VehicleDAO) UpdateVehicleOwners(rs app.RequestScope, vo *models.VehicleOwner) error {
	_, err := rs.Tx().Update("vehicle_owner", dbx.Params{
		"owner_name":  vo.OwnerName,
		"owner_email": vo.OwnerEmail,
		"owner_phone": vo.OwnerPhone},
		dbx.HashExp{"owner_id": vo.OwnerID}).Execute()
	return err
}

// VehicleOwnerExists ...
func (dao *VehicleDAO) VehicleOwnerExists(rs app.RequestScope, id uint32) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_owner WHERE owner_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// ---------------------Add / Update Fitter--------------------------------------

// CreateFitter Add Fitter
func (dao *VehicleDAO) CreateFitter(rs app.RequestScope, fd *models.FitterDetails) error {
	if exists, _ := dao.FitterExists(rs, fd.FitterID); exists == 1 {
		return dao.UpdateFitter(rs, fd)
	}
	return rs.Tx().Model(fd).Insert("FitterID", "FitterIDNo", "FittingCenterName", "FitterLocation", "FitterEmail", "FitterAddress", "FitterPhone", "FittingDate", "FitterBizRegNo")
}

// UpdateFitter update fitter
func (dao *VehicleDAO) UpdateFitter(rs app.RequestScope, fd *models.FitterDetails) error {
	_, err := rs.Tx().Update("fitter_details", dbx.Params{
		"fitting_center_name": fd.FittingCenterName,
		"fitter_location":     fd.FitterLocation,
		"fitter_email":        fd.FitterEmail,
		"fitter_address":      fd.FitterAddress,
		"fitter_phone":        fd.FitterPhone,
		"fitting_date":        fd.FittingDate,
		"fitter_biz_reg_no":   fd.FitterBizRegNo},
		dbx.HashExp{"fitting_id": fd.FitterID}).Execute()
	return err
}

// FitterExists check if fitter exists
func (dao *VehicleDAO) FitterExists(rs app.RequestScope, id uint32) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM fitter_details WHERE fitting_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// ----------------------Add / update config data----------------------

//CreateConfiguration Add configuartion details to db
func (dao *VehicleDAO) CreateConfiguration(rs app.RequestScope, cd *models.Vehicle, ownerid uint32, fitterid uint32, vehicleid uint32) error {
	a, _ := json.Marshal(cd)
	_, err := rs.Tx().Insert("vehicle_configuration", dbx.Params{
		"conf_id":           app.GenerateNewID(),
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
