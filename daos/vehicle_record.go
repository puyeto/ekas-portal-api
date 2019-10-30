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

	err := rs.Tx().Select("vehicle_id", "user_id", "owner_id", "company_id", "vehicle_string_id", "vehicle_reg_no",
		"chassis_no", "make_type", "notification_email", "notification_no", "auto_invoicing", "invoice_due_date",
		"created_on", "COALESCE(model, '')", "model_year", "COALESCE(manufacturer,'')").Model(id, &vehicleRecord)
	return &vehicleRecord, err
}

// Delete deletes an vehicleRecord with the specified ID from the database.
func (dao *VehicleRecordDAO) Delete(rs app.RequestScope, id uint32) error {
	// Delete configuration data
	_, err := rs.Tx().Delete("vehicle_configuration", dbx.HashExp{"vehicle_id": id}).Execute()
	if err != nil {
		return err
	}

	_, err = rs.Tx().Delete("vehicle_details", dbx.HashExp{"vehicle_id": id}).Execute()
	return nil
}

// Count returns the number of the vehicleRecord records in the database.
func (dao *VehicleRecordDAO) Count(rs app.RequestScope, uid int, typ string) (int, error) {
	var count int
	var err error
	if uid > 0 {
		err = rs.Tx().Select("COUNT(*)").From("vehicle_details").
			Where(dbx.HashExp{"user_id": uid}).Row(&count)
	} else {
		if typ == "ntsa" {
			err = rs.Tx().Select("COUNT(*)").From("vehicle_details").
				Where(dbx.HashExp{"send_to_ntsa": 1}).Row(&count)
		} else {
			err = rs.Tx().Select("COUNT(*)").From("vehicle_details").Row(&count)
		}
	}

	return count, err
}

// Query retrieves the vehicleRecord records with the specified offset and limit from the database.
func (dao *VehicleRecordDAO) Query(rs app.RequestScope, offset, limit int, uid int, typ string) ([]models.VehicleDetails, error) {
	vehicleRecords := []models.VehicleDetails{}
	var err error
	if uid > 0 {
		// err = rs.Tx().Select().Where(dbx.HashExp{"user_id": uid}).OrderBy("created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
		err = rs.Tx().Select("vehicle_details.vehicle_id", "vehicle_configuration.device_id", "vehicle_details.user_id", "COALESCE(company_name, '') AS company_name", "vehicle_details.vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "send_to_ntsa", "COALESCE(manufacturer, make_type) AS manufacturer", "COALESCE(model, make_type) AS model", "model_year", "vehicle_details.created_on").
			LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
			LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
			LeftJoin("vehicle_configuration", dbx.NewExp("vehicle_configuration.vehicle_string_id = vehicle_details.vehicle_string_id")).
			Where(dbx.And(dbx.HashExp{"vehicle_details.user_id": uid}, dbx.HashExp{"vehicle_configuration.status": 1}, dbx.NewExp("vehicle_configuration.device_id>0"))).
			OrderBy("vehicle_details.created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	} else {
		if typ == "ntsa" {
			err = rs.Tx().Select("vehicle_id", "vehicle_details.user_id", "COALESCE(company_name, '') AS company_name", "vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "send_to_ntsa", "COALESCE(manufacturer, make_type) AS manufacturer", "COALESCE(model, make_type) AS model", "model_year", "vehicle_details.created_on").
				LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
				LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
				Where(dbx.HashExp{"send_to_ntsa": 1}).OrderBy("vehicle_details.created_on desc").
				Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
		} else {
			err = rs.Tx().Select("vehicle_details.vehicle_id", "vehicle_configuration.device_id", "vehicle_details.user_id", "COALESCE(company_name, '') AS company_name", "vehicle_details.vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "send_to_ntsa", "COALESCE(manufacturer, make_type) AS manufacturer", "COALESCE(model, make_type) AS model", "model_year", "vehicle_details.created_on").
				LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
				LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
				LeftJoin("vehicle_configuration", dbx.NewExp("vehicle_configuration.vehicle_string_id = vehicle_details.vehicle_string_id")).
				Where(dbx.NewExp("vehicle_configuration.device_id>0")).OrderBy("vehicle_details.created_on desc").
				Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
		}
	}
	return vehicleRecords, err
}

// CreateVehicle saves a new vehicle record in the database.
func (dao *VehicleRecordDAO) CreateVehicle(rs app.RequestScope, v *models.VehicleDetails) (uint32, error) {
	err := rs.Tx().Model(v).Exclude("DeviceID", "CompanyName", "CreatedOn", "InvoiceDueDate", "AutoInvoicing", "VehicleStatus").Insert()
	return v.VehicleID, err

}

// UpdateVehicle ....
func (dao *VehicleRecordDAO) UpdateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	_, err := rs.Tx().Update("vehicle_details", dbx.Params{
		"user_id":            v.UserID,
		"vehicle_string_id":  v.VehicleStringID,
		"vehicle_reg_no":     v.VehicleRegNo,
		"chassis_no":         v.ChassisNo,
		"make_type":          v.MakeType,
		"vehicle_status":     v.VehicleStatus,
		"send_to_ntsa":       v.NTSAShow,
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

// CreateReminder saves a new reminder record in the database.
func (dao *VehicleRecordDAO) CreateReminder(rs app.RequestScope, v *models.Reminders) (uint32, error) {
	err := rs.Tx().Model(v).Insert()
	return v.ID, err

}

// CountReminders returns the number of the vehicleRecord records in the database.
func (dao *VehicleRecordDAO) CountReminders(rs app.RequestScope, uid int) (int, error) {
	var count int
	var err error
	if uid > 0 {
		err = rs.Tx().Select("COUNT(*)").From("reminders").
			Where(dbx.HashExp{"user_id": uid}).Row(&count)
	} else {
		err = rs.Tx().Select("COUNT(*)").From("reminders").Row(&count)
	}

	return count, err
}

// GetReminder retrieves the reminderRecord records with the specified offset and limit from the database.
func (dao *VehicleRecordDAO) GetReminder(rs app.RequestScope, offset, limit int, uid int) ([]models.Reminders, error) {
	vehicleRecords := []models.Reminders{}
	var err error
	if uid > 0 {
		err = rs.Tx().Select().Where(dbx.HashExp{"user_id": uid}).OrderBy("id desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	} else {
		err = rs.Tx().Select().OrderBy("id desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	}
	return vehicleRecords, err
}
