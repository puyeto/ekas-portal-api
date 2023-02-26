package daos

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

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
func (dao *TrackingServerDAO) GetTrackingServerUserLoginIDByEmail(rs app.RequestScope, email string) (interface{}, error) {
	var res struct {
		AuthUserID   uint32
		AuthUserRole int
		CompanyID    int
	}

	err := rs.Tx().Select("auth_user_id", "auth_user_role", "COALESCE(company_id, '0') AS company_id").
		From("auth_users").LeftJoin("roles", dbx.NewExp("roles.role_id = auth_users.auth_user_role")).
		LeftJoin("company_users", dbx.NewExp("company_users.user_id = auth_users.auth_user_id")).
		Where(dbx.HashExp{"auth_user_email": email}).Limit(1).One(&res)

	return res, err
}

// SaveTrackingServerLoginDetails saves a new user record in the database.
func (dao *TrackingServerDAO) SaveTrackingServerLoginDetails(rs app.RequestScope, email string, hash string, status int8, data interface{}) error {
	a, _ := json.Marshal(data)

	// insert into auth_users
	_, err := rs.Tx().Insert("auth_users", dbx.Params{
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
	err := rs.Tx().Select("COALESCE(CONCAT( first_name, ' ', last_name), '') AS full_name", "COALESCE(first_name,'') AS first_name", "COALESCE(last_name,'') AS last_name",
		"auth_user_id AS user_id", "auth_user_email AS email", "auth_user_status AS is_verified", "auth_user_role AS role", "role_name",
		"COALESCE(company_users.company_id, 0) AS company_id", "COALESCE(company_name, '') AS company_name", "enable_gps_configuration", "mpesa_renewal",
		"enable_failsafe_configuration", "sacco_id", "auth_user_password", "auth_user_salt").
		From("auth_users").
		LeftJoin("roles", dbx.NewExp("roles.role_id = auth_users.auth_user_role")).
		LeftJoin("company_users", dbx.NewExp("company_users.user_id = auth_users.auth_user_id")).
		LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).
		Where(dbx.HashExp{"auth_user_email": email}).One(&usr)

	return usr, err
}

// GetCompanyDetailsByEmail ...
func (dao *TrackingServerDAO) GetCompanyDetailsByEmail(rs app.RequestScope, email string) (models.Companies, error) {
	com := models.Companies{}
	err := rs.Tx().Select("companies.*").
		From("companies").
		LeftJoin("company_users", dbx.NewExp("company_users.company_id = companies.company_id")).
		LeftJoin("auth_users", dbx.NewExp("auth_users.auth_user_id = company_users.user_id")).
		Where(dbx.HashExp{"auth_user_email": email}).One(&com)

	return com, err
}

// GetUserByUserHash reads the user with the specified email from the database.
func (dao *TrackingServerDAO) GetUserByUserHash(rs app.RequestScope, userhash string) (models.AdminUserDetails, error) {
	usr := models.AdminUserDetails{}
	err := rs.Tx().Select("COALESCE(first_name, '') AS user_first_name", "COALESCE(last_name, '') AS user_last_name", "auth_user_id AS user_id", "auth_user_email AS email").
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

// Get reads the sacco with the specified ID from the database.
func (dao *TrackingServerDAO) GetSaccoName(rs app.RequestScope, id int) (string, error) {
	var sacconame string
	err := rs.Tx().NewQuery("SELECT name FROM saccos WHERE id='" + strconv.Itoa(id) + "' LIMIT 1").Row(&sacconame)
	return sacconame, err
}

// ResetPassword Reset admin password.
func (dao *TrackingServerDAO) ResetPassword(rs app.RequestScope, password string, userid int32) error {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	salt := fmt.Sprintf("%x", h.Sum(nil))
	hashpassword := app.CalculatePassHash(password, salt)

	_, err := rs.Tx().Update("auth_users", dbx.Params{
		"auth_user_password": hashpassword,
		"auth_user_salt":     salt},
		dbx.HashExp{"auth_user_id": userid}).Execute()

	return err
}

func (dao *TrackingServerDAO) ResetPassword2(rs app.RequestScope, password string, userid int32) (hashpassword, salt string, err error) {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	salt = fmt.Sprintf("%x", h.Sum(nil))
	hashpassword = app.CalculatePassHash(password, salt)

	_, err = rs.Tx().Update("auth_users", dbx.Params{
		"auth_user_password": hashpassword,
		"auth_user_salt":     salt},
		dbx.HashExp{"auth_user_id": userid}).Execute()

	return hashpassword, salt, err
}
