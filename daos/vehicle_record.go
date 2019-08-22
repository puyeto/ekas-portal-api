package daos

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// VehicleRecordDAO persists vehicleRecord data in database
type VehicleRecordDAO struct{}

// NewVehicleRecordDAO creates a new VehicleRecordDAO
func NewVehicleRecordDAO() *VehicleRecordDAO {
	return &VehicleRecordDAO{}
}

// Get reads the vehicleRecord with the specified ID from the database.
func (dao *VehicleRecordDAO) Get(rs app.RequestScope, id uint32) (*models.VehicleDetails, error) {
	var vehicleRecord models.VehicleDetails
	err := rs.Tx().Select().Model(id, &vehicleRecord)
	return &vehicleRecord, err
}

// Update saves the changes to an vehicleRecord in the database.
func (dao *VehicleRecordDAO) Update(rs app.RequestScope, id uint32, vehicleRecord *models.VehicleDetails) error {
	if _, err := dao.Get(rs, id); err != nil {
		return err
	}
	vehicleRecord.VehicleID = id
	return rs.Tx().Model(vehicleRecord).Exclude("Id").Update()
}

// Delete deletes an vehicleRecord with the specified ID from the database.
func (dao *VehicleRecordDAO) Delete(rs app.RequestScope, id uint32) error {
	vehicleRecord, err := dao.Get(rs, id)
	if err != nil {
		return err
	}
	return rs.Tx().Model(vehicleRecord).Delete()
}

// Count returns the number of the vehicleRecord records in the database.
func (dao *VehicleRecordDAO) Count(rs app.RequestScope, uid int) (int, error) {
	var count int
	var err error
	if uid > 0 {
		err = rs.Tx().Select("COUNT(*)").From("vehicle_details").
			Where(dbx.HashExp{"user_id": uid}).Row(&count)
	} else {
		err = rs.Tx().Select("COUNT(*)").From("vehicle_details").Row(&count)
	}

	return count, err
}

// Query retrieves the vehicleRecord records with the specified offset and limit from the database.
func (dao *VehicleRecordDAO) Query(rs app.RequestScope, offset, limit int, uid int) ([]models.VehicleDetails, error) {
	vehicleRecords := []models.VehicleDetails{}
	var err error
	if uid > 0 {
		err = rs.Tx().Select().Where(dbx.HashExp{"user_id": uid}).OrderBy("created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	} else {
		err = rs.Tx().Select("vehicle_id", "vehicle_details.user_id", "COALESCE(company_name, '') AS company_name", "vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "vehicle_details.created_on").
			LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
			LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
			OrderBy("vehicle_details.created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	}
	return vehicleRecords, err
}

// CreateVehicle saves a new vehicle record in the database.
func (dao *VehicleRecordDAO) CreateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	return rs.Tx().Model(v).Insert("VehicleID", "UserID", "VehicleStringID", "VehicleRegNo", "ChassisNo", "MakeType", "NotificationEmail", "NotificationNO")
}

// UpdateVehicle ....
func (dao *VehicleRecordDAO) UpdateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
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
func (dao *VehicleRecordDAO) VehicleExists(rs app.RequestScope, id uint32) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_details WHERE vehicle_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}
