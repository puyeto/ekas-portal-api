package daos

import (
	"encoding/json"
	"fmt"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// TrackingServerDAO persists trackingServer data in database
type TrackingServerDAO struct{}

// NewTrackingServerDAO creates a new TrackingServerDAO
func NewTrackingServerDAO() *TrackingServerDAO {
	return &TrackingServerDAO{}
}

// GetTrackingServerUserLoginIDByEmail ...
func (dao *TrackingServerDAO) GetTrackingServerUserLoginIDByEmail(rs app.RequestScope, email string) (uint32, int, int, error) {
	var res struct {
		AuthUserID   uint32
		AuthUserRole int
		CompanyID    int
	}

	err := rs.Tx().Select("auth_user_id", "auth_user_role", "COALESCE(company_id, '0') AS company_id").
		From("auth_users").LeftJoin("roles", dbx.NewExp("roles.role_id = auth_users.auth_user_role")).
		LeftJoin("company_users", dbx.NewExp("company_users.user_id = auth_users.auth_user_id")).
		Where(dbx.HashExp{"auth_user_email": email}).Limit(1).One(&res)

	// q := rs.Tx().NewQuery("SELECT auth_user_id, auth_user_role FROM auth_users WHERE auth_user_email='" + email + "' LIMIT 1")
	// err := q.One(&res)
	fmt.Println(res)
	return res.AuthUserID, res.AuthUserRole, res.CompanyID, err
}

// SaveTrackingServerLoginDetails saves a new user record in the database.
func (dao *TrackingServerDAO) SaveTrackingServerLoginDetails(rs app.RequestScope, id uint32, email string, hash string, status int8, data interface{}) error {
	//return rs.Tx().Model(artist).Insert()
	a, _ := json.Marshal(data)

	// insert into auth_users
	_, err := rs.Tx().Insert("auth_users", dbx.Params{
		"auth_user_id":     id,
		"auth_user_email":  string(email),
		"auth_user_hash":   string(hash),
		"auth_user_status": status,
		"auth_user_data":   string(a),
	}).Execute()
	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail reads the user with the specified email from the database.
func (dao *TrackingServerDAO) GetUserByEmail(rs app.RequestScope, email string) (models.AdminUserDetails, error) {
	usr := models.AdminUserDetails{}
	err := rs.Tx().Select("COALESCE(first_name,'')", "COALESCE(last_name,'')", "auth_user_id AS user_id", "auth_user_email AS email", "auth_user_hash AS token", "auth_user_status AS is_verified", "auth_user_role AS role_id", "role_name").
		From("auth_users").
		LeftJoin("roles", dbx.NewExp("roles.role_id = auth_users.auth_user_role")).
		Where(dbx.HashExp{"auth_user_email": email}).One(&usr)

	return usr, err
}

// GetUserByUserHash reads the user with the specified email from the database.
func (dao *TrackingServerDAO) GetUserByUserHash(rs app.RequestScope, userhash string) (models.AdminUserDetails, error) {
	usr := models.AdminUserDetails{}
	err := rs.Tx().Select("COALESCE(first_name, '')", "COALESCE(last_name, '')", "auth_user_id AS user_id", "auth_user_email AS email").
		From("auth_users").
		Where(dbx.HashExp{"auth_user_hash": userhash}).One(&usr)

	return usr, err
}

// CreateLoginSession creates a new one-time-use login token
func (dao *TrackingServerDAO) CreateLoginSession(rs app.RequestScope, ls *models.UserLoginSessions) error {
	return rs.Tx().Model(ls).Exclude().Insert()
}

// TrackingServerUserEmailExists check if a tracker server auth user exists
func (dao *TrackingServerDAO) TrackingServerUserEmailExists(rs app.RequestScope, email string) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM auth_users WHERE auth_user_email='" + email + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// QueryVehicelsFromPortal retrieves the vehicleRecord records with the specified offset and limit from the database.
func (dao *TrackingServerDAO) QueryVehicelsFromPortal(rs app.RequestScope, offset, limit int, uid int) ([]models.VehicleDetails, error) {
	vehicleRecords := []models.VehicleDetails{}
	var err error
	if uid > 0 {
		// err = rs.Tx().Select().Where(dbx.HashExp{"user_id": uid}).OrderBy("created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
		err = rs.Tx().Select("vehicle_details.vehicle_id", "vehicle_details.user_id", "COALESCE(device_id, 0) AS device_id", "COALESCE(company_name, '') AS company_name", "vehicle_details.vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "COALESCE(manufacturer, make_type) AS manufacturer", "COALESCE(model, make_type) AS model", "model_year", "vehicle_details.created_on").
			LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
			LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
			LeftJoin("vehicle_configuration", dbx.And(dbx.NewExp("vehicle_configuration.vehicle_id = vehicle_details.vehicle_id"), dbx.NewExp("vehicle_configuration.status=1"))).
			Where(dbx.HashExp{"vehicle_details.user_id": uid}).
			OrderBy("vehicle_details.created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	} else {
		err = rs.Tx().Select("vehicle_details.vehicle_id", "vehicle_details.user_id", "COALESCE(device_id, 0) AS device_id", "COALESCE(company_name, '') AS company_name", "vehicle_details.vehicle_string_id", "vehicle_reg_no", "chassis_no", "make_type", "notification_email", "notification_no", "vehicle_status", "COALESCE(manufacturer, make_type) AS manufacturer", "COALESCE(model, make_type) AS model", "model_year", "vehicle_details.created_on").
			LeftJoin("company_users", dbx.NewExp("company_users.user_id = vehicle_details.user_id")).
			LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
			LeftJoin("vehicle_configuration", dbx.And(dbx.NewExp("vehicle_configuration.vehicle_id = vehicle_details.vehicle_id"), dbx.NewExp("vehicle_configuration.status=1"))).
			OrderBy("vehicle_details.created_on desc").Offset(int64(offset)).Limit(int64(limit)).All(&vehicleRecords)
	}
	return vehicleRecords, err
}
