package daos

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// UserDAO persists user data in database
type UserDAO struct{}

// NewUserDAO creates a new UserDAO
func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

// Count returns the number of the company records in the database.
func (dao *UserDAO) Count(rs app.RequestScope, cid int) (int, error) {
	var count int
	q := rs.Tx().Select("COUNT(*)").From("auth_users")
	if cid > 0 {
		q.LeftJoin("company_users", dbx.NewExp("company_users.user_id = auth_users.auth_user_id")).
			Where(dbx.HashExp{"company_users.company_id": cid})
	}
	err := q.Row(&count)
	return count, err
}

// Query retrieves the company records with the specified offset and limit from the database.
func (dao *UserDAO) Query(rs app.RequestScope, offset, limit, cid int) ([]models.AuthUsers, error) {
	users := []models.AuthUsers{}
	q := rs.Tx().Select("auth_user_id", "auth_user_email", "auth_user_status", "auth_user_role", "sacco_id", "role_name", "COALESCE(CONCAT( first_name, ' ', last_name), '') AS full_name", "COALESCE( first_name, '') AS first_name", "COALESCE(last_name, '') AS last_name, COALESCE(companies.company_id, 0) AS company_id, COALESCE(company_name, '') AS company_name").
		From("auth_users").LeftJoin("roles", dbx.NewExp("roles.role_id = auth_users.auth_user_role")).
		LeftJoin("company_users", dbx.NewExp("company_users.user_id = auth_users.auth_user_id")).
		LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id"))
	if cid > 0 {
		q.Where(dbx.HashExp{"company_users.company_id": cid})
	}
	err := q.OrderBy("auth_user_id ASC").Offset(int64(offset)).Limit(int64(limit)).All(&users)
	return users, err
}

// GetUser reads the full user details with the specified ID from the database.
func (dao *UserDAO) GetUser(rs app.RequestScope, id uint32) (models.AuthUsers, error) {
	usr := models.AuthUsers{}
	err := rs.Tx().Select("auth_user_id", "auth_user_email", "auth_user_status", "sacco_id", "auth_user_role", "role_name", "COALESCE( first_name, '') AS first_name", "COALESCE(last_name, '') AS last_name, COALESCE(companies.company_id, 0) AS company_id, COALESCE(company_name, '') AS company_name").
		From("auth_users").LeftJoin("roles", dbx.NewExp("roles.role_id = auth_users.auth_user_role")).
		LeftJoin("company_users", dbx.NewExp("company_users.user_id = auth_users.auth_user_id")).
		LeftJoin("companies", dbx.NewExp("companies.company_id = company_users.company_id")).Model(id, &usr)
	if err != nil {
		return usr, err
	}

	return usr, err
}

// Update ....
func (dao *UserDAO) Update(rs app.RequestScope, a *models.AuthUsers) error {
	_, err := rs.Tx().Update("auth_users", dbx.Params{
		"auth_user_email":  a.Email,
		"auth_user_role":   a.RoleID,
		"auth_user_status": a.Status,
		"first_name":       a.FirstName,
		"last_name":        a.LastName,
		"sacco_id":         a.SaccoID},
		dbx.HashExp{"auth_user_id": a.UserID}).Execute()
	return err
}

// ResetPassword Reset admin password.
func (dao *UserDAO) ResetPassword(rs app.RequestScope, m *models.ResetPassword) error {
	_, err := rs.Tx().Update("admin_user_details", dbx.Params{
		"password": m.Password,
		"salt":     m.Salt},
		dbx.HashExp{"user_id": m.UserID}).Execute()

	return err
}

// QueryDepartments retrieves the company records with the specified offset and limit from the database.
func (dao *UserDAO) QueryDepartments(rs app.RequestScope) ([]models.Departments, error) {
	dep := []models.Departments{}
	err := rs.Tx().Select("id", "name", "description").From("departments").OrderBy("id ASC").All(&dep)
	return dep, err
}

// CreateCompanyUser create user relationship to company.
func (dao *UserDAO) CreateCompanyUser(rs app.RequestScope, companyid int32, userid uint32) error {
	_, err := rs.Tx().Insert("company_users", dbx.Params{
		"user_id":    userid,
		"company_id": companyid}).Execute()
	return err
}

// UpdateCompanyUser Update user relationship to company.
func (dao *UserDAO) UpdateCompanyUser(rs app.RequestScope, companyid int32, userid uint32) error {
	_, err := rs.Tx().Update("company_users", dbx.Params{
		"company_id": companyid},
		dbx.HashExp{"user_id": userid}).Execute()

	return err
}

// IfCompanyUserExists check if company and user exists (company_users)
func (dao *UserDAO) IfCompanyUserExists(rs app.RequestScope, companyid int32, userid uint32) (int, error) {
	var exists int
	// q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM company_users WHERE user_id='" + strconv.Itoa(int(userid)) + "' AND company_id='" + strconv.Itoa(int(companyid)) + "' LIMIT 1) AS exist")
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM company_users WHERE user_id='" + strconv.Itoa(int(userid)) + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// // GetUserByEmail reads the user with the specified email from the database.
// func (dao *UserDAO) GetUserByEmail(rs app.RequestScope, email string) (*models.UserDetails, error) {
// 	var usr models.UserDetails
// 	err := rs.Tx().Select("first_name", "last_name", "user_id", "mobile_number", "email", "password", "username", "is_verified", "salt").
// 		Where(dbx.HashExp{"email": email}).One(&usr)

// 	return &usr, err
// }

// GetUserByEmail reads the user with the specified email from the database.
func (dao *UserDAO) GetUserByEmail(rs app.RequestScope, email string) (*models.AdminUserDetails, error) {
	usr := &models.AdminUserDetails{}
	err := rs.Tx().Select("first_name", "last_name", "aud.user_id", "mobile_number", "email", "password", "username", "is_verified", "salt", "aur.role_id", "role_name").
		From("admin_user_details AS aud").
		LeftJoin("admin_user_roles AS aur", dbx.NewExp("aur.user_id = aud.user_id")).
		LeftJoin("roles", dbx.NewExp("roles.role_id = aur.role_id")).
		Where(dbx.HashExp{"email": email}).One(&usr)

	if err != nil {
		return nil, err
	}

	return usr, nil
}

// Register saves a new user record in the database.
// The User.ID field will be populated with an automatically generated ID upon successful saving.
func (dao *UserDAO) Register(rs app.RequestScope, usr *models.AdminUserDetails) error {
	return rs.Tx().Model(usr).Insert("UserID", "FirstName", "LastName", "Email", "Password", "Username", "DOB", "MobileNumber", "Salt", "VerificationCode")
}

// SubmitUserRole submit user role
func (dao *UserDAO) SubmitUserRole(rs app.RequestScope, ur *models.AdminUserRoles) error {
	if ur.RoleID == 0 {
		ur.RoleID = 10001
	}

	return rs.Tx().Model(ur).Exclude("RoleName").Insert()
}

// Delete deletes user with the specified ID from the database.
func (dao *UserDAO) Delete(rs app.RequestScope, id int32) error {
	_, err := rs.Tx().Delete("user_details", dbx.HashExp{"user_id": id}).Execute()
	return err
}

// CreateNewEmailVerification - Create a new user
func (dao *UserDAO) CreateNewEmailVerification(rs app.RequestScope, con *models.ConfirmationEmailDetails) error {

	if err := con.VerifyConfirmationEmail(); err != nil {
		return err
	}

	_, err := rs.Tx().Update("user_details", dbx.Params{
		"verification_code": con.VerificationCode},
		dbx.HashExp{"user_id": con.UserID}).Execute()

	return err
}

// IsEmailExists ...
func (dao *UserDAO) IsEmailExists(rs app.RequestScope, email string) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM admin_user_details WHERE email='" + email + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}

// CreateLoginSession creates a new one-time-use login token
func (dao *UserDAO) CreateLoginSession(rs app.RequestScope, ls *models.UserLoginSessions) error {
	return rs.Tx().Model(ls).Exclude().Insert()
}
