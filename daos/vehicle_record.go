package daos

import (
	"errors"
	"strconv"
	"strings"
	"time"

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
		err = rs.Tx().Select("vehicle_details.vehicle_id", "vehicle_configuration.device_id", "vehicle_details.user_id", "COALESCE(company_name, '') AS company_name", "vehicle_details.vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "send_to_ntsa", "COALESCE(manufacturer, make_type) AS manufacturer", "COALESCE(model, make_type) AS model", "model_year", "vehicle_details.created_on", "last_seen", "COALESCE(renewal_date, vehicle_details.created_on) AS renewal_date", "renew").
			LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
			LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
			LeftJoin("vehicle_configuration", dbx.NewExp("vehicle_configuration.vehicle_string_id = vehicle_details.vehicle_string_id")).
			Where(dbx.And(dbx.HashExp{"vehicle_details.user_id": uid}, dbx.HashExp{"vehicle_details.vehicle_status": 1}, dbx.HashExp{"vehicle_configuration.status": 1}, dbx.NewExp("vehicle_configuration.device_id>0"))).
			OrderBy("vehicle_details.created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	} else {
		if typ == "ntsa" {
			err = rs.Tx().Select("vehicle_id", "vehicle_details.user_id", "COALESCE(company_name, '') AS company_name", "vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "send_to_ntsa", "COALESCE(manufacturer, make_type) AS manufacturer", "COALESCE(model, make_type) AS model", "model_year", "vehicle_details.created_on", "last_seen", "COALESCE(renewal_date, vehicle_details.created_on) AS renewal_date", "renew").
				LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
				LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
				Where(dbx.HashExp{"send_to_ntsa": 1}).OrderBy("vehicle_details.created_on desc").
				Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
		} else {
			err = rs.Tx().Select("vehicle_details.vehicle_id", "vehicle_configuration.device_id", "vehicle_details.user_id", "COALESCE(company_name, '') AS company_name", "vehicle_details.vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "send_to_ntsa", "COALESCE(manufacturer, make_type) AS manufacturer", "COALESCE(model, make_type) AS model", "model_year", "vehicle_details.created_on", "last_seen", "COALESCE(renewal_date, vehicle_details.created_on) AS renewal_date", "renew").
				LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
				LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
				LeftJoin("vehicle_configuration", dbx.NewExp("vehicle_configuration.vehicle_string_id = vehicle_details.vehicle_string_id")).
				Where(dbx.And(dbx.NewExp("vehicle_configuration.device_id>0"), dbx.HashExp{"vehicle_details.vehicle_status": 1})).OrderBy("vehicle_details.created_on desc").
				Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
		}
	}
	return vehicleRecords, err
}

// QueryFilter retrieves the filtered vehicleRecord records with the specified offset and limit from the database.
func (dao *VehicleRecordDAO) QueryFilter(rs app.RequestScope, offset, limit int, model *models.FilterVehicles) ([]models.VehicleDetails, error) {
	vehicleRecords := []models.VehicleDetails{}
	q := rs.Tx().Select("vehicle_details.vehicle_id", "vehicle_configuration.device_id", "vehicle_details.user_id",
		"COALESCE(company_name, '') AS company_name", "vehicle_details.vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "JSON_VALUE(data, '$.device_detail.owner_name') AS vehicle_owner",
		"notification_email", "notification_no", "vehicle_status", "send_to_ntsa", "JSON_VALUE(data, '$.device_detail.serial_no') AS limiter_serial", "JSON_VALUE(data, '$.device_detail.certificate') AS certificate",
		"COALESCE(model, make_type) AS model", "limiter_type", "JSON_VALUE(data, '$.device_detail.owner_phone_number') AS vehicle_owner_tel", "JSON_VALUE(data, '$.device_detail.agent_location') AS fitting_location", "vehicle_details.created_on", "last_seen").
		LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
		LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
		LeftJoin("vehicle_configuration", dbx.NewExp("vehicle_configuration.vehicle_string_id = vehicle_details.vehicle_string_id"))

	if model.MinTimeStamp != "" && model.MaxTimeStamp != "" && model.FilterNTSA != 2 && model.FilterStatus != 2 {
		q.Where(dbx.And(dbx.NewExp("vehicle_configuration.device_id>0"), dbx.Between("vehicle_details.created_on", model.MinTimeStamp, model.MaxTimeStamp), dbx.HashExp{"send_to_ntsa": model.FilterNTSA}, dbx.HashExp{"vehicle_status": model.FilterStatus}))
	} else if model.MinTimeStamp != "" && model.MaxTimeStamp != "" && model.FilterNTSA != 2 {
		q.Where(dbx.And(dbx.NewExp("vehicle_configuration.device_id>0"), dbx.Between("vehicle_details.created_on", model.MinTimeStamp, model.MaxTimeStamp), dbx.HashExp{"send_to_ntsa": model.FilterNTSA}))
	} else if model.MinTimeStamp != "" && model.MaxTimeStamp != "" && model.FilterStatus != 2 {
		q.Where(dbx.And(dbx.NewExp("vehicle_configuration.device_id>0"), dbx.Between("vehicle_details.created_on", model.MinTimeStamp, model.MaxTimeStamp), dbx.HashExp{"vehicle_status": model.FilterStatus}))
	} else if model.FilterNTSA != 2 && model.FilterStatus != 2 {
		q.Where(dbx.And(dbx.NewExp("vehicle_configuration.device_id>0"), dbx.HashExp{"send_to_ntsa": model.FilterNTSA}, dbx.HashExp{"vehicle_status": model.FilterStatus}))
	} else if model.MinTimeStamp != "" && model.MaxTimeStamp != "" {
		q.Where(dbx.And(dbx.NewExp("vehicle_configuration.device_id>0"), dbx.Between("vehicle_details.created_on", model.MinTimeStamp, model.MaxTimeStamp)))
	} else if model.FilterNTSA != 2 {
		q.Where(dbx.And(dbx.NewExp("vehicle_configuration.device_id>0"), dbx.HashExp{"send_to_ntsa": model.FilterNTSA}))
	} else if model.FilterStatus != 2 {
		q.Where(dbx.And(dbx.NewExp("vehicle_configuration.device_id>0"), dbx.HashExp{"vehicle_status": model.FilterStatus}))
	}

	err := q.OrderBy("vehicle_details.created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	return vehicleRecords, err
}

// CountFilter returns the number of the filtered vehicleRecord records in the database.
func (dao *VehicleRecordDAO) CountFilter(rs app.RequestScope, model *models.FilterVehicles) (int, error) {
	var count int
	var q = rs.Tx().Select("COUNT(*)").From("vehicle_details")
	if model.MinTimeStamp != "" && model.MaxTimeStamp != "" && model.FilterNTSA != 2 && model.FilterStatus != 2 {
		q.Where(dbx.And(dbx.Between("created_on", model.MinTimeStamp, model.MaxTimeStamp), dbx.HashExp{"send_to_ntsa": model.FilterNTSA}, dbx.HashExp{"vehicle_status": model.FilterStatus}))
	} else if model.MinTimeStamp != "" && model.MaxTimeStamp != "" && model.FilterNTSA != 2 {
		q.Where(dbx.And(dbx.Between("created_on", model.MinTimeStamp, model.MaxTimeStamp), dbx.HashExp{"send_to_ntsa": model.FilterNTSA}))
	} else if model.MinTimeStamp != "" && model.MaxTimeStamp != "" && model.FilterStatus != 2 {
		q.Where(dbx.And(dbx.Between("created_on", model.MinTimeStamp, model.MaxTimeStamp), dbx.HashExp{"vehicle_status": model.FilterStatus}))
	} else if model.FilterNTSA != 2 && model.FilterStatus != 2 {
		q.Where(dbx.And(dbx.HashExp{"send_to_ntsa": model.FilterNTSA}, dbx.HashExp{"vehicle_status": model.FilterStatus}))
	} else if model.MinTimeStamp != "" && model.MaxTimeStamp != "" {
		q.Where(dbx.Between("created_on", model.MinTimeStamp, model.MaxTimeStamp))
	} else if model.FilterNTSA != 2 {
		q.Where(dbx.HashExp{"send_to_ntsa": model.FilterNTSA})
	} else if model.FilterStatus != 2 {
		q.Where(dbx.HashExp{"vehicle_status": model.FilterStatus})
	}

	err := q.Row(&count)

	return count, err
}

// CreateVehicle saves a new vehicle record in the database.
func (dao *VehicleRecordDAO) CreateVehicle(rs app.RequestScope, v *models.VehicleDetails) (uint32, error) {
	err := rs.Tx().Model(v).Exclude("DeviceID", "CompanyName", "CreatedOn", "InvoiceDueDate", "AutoInvoicing", "VehicleStatus").Insert()
	return v.VehicleID, err

}

// UpdateVehicle ....
func (dao *VehicleRecordDAO) UpdateVehicle(rs app.RequestScope, v *models.VehicleDetails) error {
	if _, err := rs.Tx().Update("vehicle_details", dbx.Params{
		"user_id":            v.UserID,
		"vehicle_string_id":  v.VehicleStringID,
		"vehicle_reg_no":     strings.ToUpper(v.VehicleRegNo),
		"chassis_no":         strings.ToUpper(v.ChassisNo),
		"make_type":          strings.ToUpper(v.MakeType),
		"vehicle_status":     v.VehicleStatus,
		"send_to_ntsa":       v.NTSAShow,
		"notification_email": v.NotificationEmail,
		"notification_no":    v.NotificationNO},
		dbx.HashExp{"vehicle_id": v.VehicleID}).Execute(); err != nil {
		return err
	}

	// update configuration details
	query := "UPDATE vehicle_configuration SET vehicle_string_id = '" + v.VehicleStringID
	query += "', data = JSON_SET(DATA, '$.device_detail.registration_no', '" + v.VehicleRegNo + "', '$.device_detail.chasis_no', '" + v.ChassisNo + "', '$.device_detail.make_type', '" + v.MakeType + "')"
	query += " WHERE vehicle_id = " + strconv.Itoa(int(v.VehicleID))
	if _, err := rs.Tx().NewQuery(query).Execute(); err != nil {
		return err
	}

	return nil
}

// VehicleExists ...
func (dao *VehicleRecordDAO) VehicleExists(rs app.RequestScope, id uint32) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_details WHERE vehicle_id='" + strconv.Itoa(int(id)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// RenewVehicle ...
func (dao *VehicleRecordDAO) RenewVehicle(rs app.RequestScope, m *models.VehicleRenewals) (uint32, error) {
	m.Status = 1
	m.CreatedOn = time.Now()
	m.RenewalDate = m.RenewalDate.AddDate(1, 0, 0)
	m.ExpiryDate = m.RenewalDate.AddDate(1, 0, -1)

	// check if cert has been renewed
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM vehicle_renewals WHERE certificate_no='" + m.CertificateNo + "' LIMIT 1) AS exist")
	if err := q.Row(&exists); err != nil {
		return m.ID, err
	}

	if exists == 1 {
		return m.ID, errors.New("Certificate has been renewed")
	}

	// Save renewal details
	if err := rs.Tx().Model(m).Exclude("VehicleRegNo", "DeviceSerialNo").Insert(); err != nil {
		return m.ID, err
	}

	// update vehicle details
	if _, err := rs.Tx().Update("vehicle_details", dbx.Params{
		"renew":        m.Status,
		"renewal_date": m.ExpiryDate},
		dbx.HashExp{"vehicle_id": m.VehicleID}).Execute(); err != nil {
		rs.Rollback()
		return m.ID, err
	}

	return m.ID, nil
}

// CreateReminder saves a new reminder record in the database.
func (dao *VehicleRecordDAO) CreateReminder(rs app.RequestScope, v *models.Reminders) (uint32, error) {
	err := rs.Tx().Model(v).Insert()
	return v.ID, err

}

// ListVehicleRenewals retrieves the renewal records with the specified offset and limit from the database.
func (dao *VehicleRecordDAO) ListVehicleRenewals(rs app.RequestScope, offset, limit int) ([]models.VehicleRenewals, error) {
	r := []models.VehicleRenewals{}
	err := rs.Tx().Select("id", "serial_no AS device_serial_no", "vehicle_reg_no", "vr.vehicle_id", "vr.vehicle_string_id", "vr.status", "added_by", "vr.renewal_date", "vr.expiry_date", "vr.renewal_code", "vr.created_on").
		From("vehicle_renewals AS vr").
		LeftJoin("vehicle_details", dbx.NewExp("vehicle_details.vehicle_id = vr.vehicle_id")).
		LeftJoin("vehicle_configuration AS vc", dbx.NewExp("vc.vehicle_id = vr.vehicle_id")).
		OrderBy("id DESC").Offset(int64(offset)).Limit(int64(limit)).All(&r)
	return r, err
}

// CountRenewals returns the number of the renewals records in the database.
func (dao *VehicleRecordDAO) CountRenewals(rs app.RequestScope) (int, error) {
	var count int
	err := rs.Tx().Select("COUNT(*)").From("vehicle_renewals").Row(&count)

	return count, err
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
