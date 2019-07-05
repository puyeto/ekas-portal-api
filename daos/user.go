package daos

import (
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

// GetUser reads the full user details with the specified ID from the database.
func (dao *UserDAO) GetUser(rs app.RequestScope, id int32) (*models.ListUserDetails, error) {
	usr := &models.ListUserDetails{}
	err := rs.Tx().Select("first_name", "last_name", "user_id", "mobile_number", "email", "username", "is_verified").
		From("user_details").Model(id, &usr.UserDetails)
	if err != nil {
		return nil, err
	}

	var rls []models.AdminUserRoles
	//rs.Tx().Select().Where(dbx.HashExp{"user_id": id}).All(&rls)
	rs.Tx().Select().From("user_roles").
		LeftJoin("roles", dbx.NewExp("roles.role_id = user_roles.role_id")).
		Where(dbx.HashExp{"user_id": id}).All(&rls)
	usr.Roles = rls

	return usr, err
}

// // GetUserByEmail reads the user with the specified email from the database.
// func (dao *UserDAO) GetUserByEmail(rs app.RequestScope, email string) (*models.UserDetails, error) {
// 	var usr models.UserDetails
// 	err := rs.Tx().Select("first_name", "last_name", "user_id", "mobile_number", "email", "password", "username", "is_verified", "salt").
// 		Where(dbx.HashExp{"email": email}).One(&usr)

// 	return &usr, err
// }

// GetUserByEmail reads the user with the specified email from the database.
func (dao *UserDAO) GetUserByEmail(rs app.RequestScope, email string) (*models.ListUserDetails, error) {
	usr := &models.ListUserDetails{}
	err := rs.Tx().Select("first_name", "last_name", "user_id", "mobile_number", "email", "password", "username", "is_verified", "salt").
		Where(dbx.HashExp{"email": email}).One(&usr.UserDetails)

	if err != nil {
		return nil, err
	}

	var rls []models.AdminUserRoles
	//rs.Tx().Select().Where(dbx.HashExp{"user_id": usr.UserDetails.UserID}).All(&rls)
	rs.Tx().Select().From("user_roles").
		LeftJoin("roles", dbx.NewExp("roles.role_id = user_roles.role_id")).
		Where(dbx.HashExp{"user_id": usr.UserDetails.UserID}).All(&rls)
	usr.Roles = rls

	return usr, err
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
