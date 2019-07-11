package services

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	"github.com/go-ozzo/ozzo-routing/auth"
	uuid "github.com/nu7hatch/gouuid"
)

// userDAO specifies the interface of the user DAO needed by userService.
type userDAO interface {
	// GetUser returns the user with the specified user ID.
	GetUser(rs app.RequestScope, id int32) (*models.ListUserDetails, error)
	// GetUserByEmail returns the user with the specified user email.
	// GetUserByEmail(rs app.RequestScope, email string) (*models.UserDetails, error)
	GetUserByEmail(rs app.RequestScope, email string) (*models.AdminUserDetails, error)
	// Register saves a new user in the storage.
	Register(rs app.RequestScope, usr *models.AdminUserDetails) error
	IsEmailExists(rs app.RequestScope, email string) (int, error)
	// SubmitUserRole saves a new user role.
	SubmitUserRole(rs app.RequestScope, usr *models.AdminUserRoles) error
	Delete(rs app.RequestScope, id int32) error
	CreateNewEmailVerification(rs app.RequestScope, con *models.ConfirmationEmailDetails) error
	CreateLoginSession(rs app.RequestScope, ls *models.UserLoginSessions) error
}

// UserService provides services related with users.
type UserService struct {
	dao userDAO
}

// NewUserService creates a new userService with the given user DAO.
func NewUserService(dao userDAO) *UserService {
	return &UserService{dao}
}

// New new user salt
func New() models.AdminUserDetails {
	u := models.AdminUserDetails{}
	u.NewSalt()
	return u
}

// GetUser returns the user with the specified the user ID.
func (u *UserService) GetUser(rs app.RequestScope, id int32) (*models.ListUserDetails, error) {
	return u.dao.GetUser(rs, id)
}

// GetUserByEmail returns the user with the specified the user email.
func (u *UserService) GetUserByEmail(rs app.RequestScope, email string) (*models.AdminUserDetails, error) {
	return u.dao.GetUserByEmail(rs, email)
}

//Login a user
func (u *UserService) Login(rs app.RequestScope, c *models.Credential) (*models.AdminUserDetails, error) {
	if err := c.ValidateCredential(); err != nil {
		return nil, err
	}

	res, err := u.dao.GetUserByEmail(rs, c.Email)
	if err != nil {
		return nil, err
	}

	if &res == nil {
		return nil, errors.New("no user found")
	}

	if res.Password != app.CalculatePassHash(c.Password, res.Salt) {
		return nil, errors.New("invalid credential")
	}

	if res.IsVerified != 1 {
		return nil, errors.New("Account not verified")
	}

	reset(res)

	token, err := auth.NewJWT(jwt.MapClaims{
		"id":  strconv.Itoa(int(res.UserID)),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}, app.Config.JWTSigningKey)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	res.Token = token

	u.storeLoginSession(rs, res)

	return res, nil
}

func (u *UserService) storeLoginSession(rs app.RequestScope, ud *models.AdminUserDetails) error {
	r := &http.Request{}
	log.Println(r)
	loginSession := models.UserLoginSessions{
		SessionID: app.GenerateNewID(),
		UserID:    ud.UserID,
		UserAgent: r.UserAgent(),
		IP:        models.GetRemoteIP(r),
		Token:     ud.Token,
	}

	return u.dao.CreateLoginSession(rs, &loginSession)
}

func reset(u *models.AdminUserDetails) {
	// reset password and salt
	u.Password = ""
	u.Salt = ""
}

// Register creates a new user.
func (u *UserService) Register(rs app.RequestScope, model *models.AdminUserDetails) (int32, error) {
	model.UserID = 0
	if err := model.ValidateNewUser(); err != nil {
		return 0, err
	}

	//cehck if email exists
	exists, err := u.dao.IsEmailExists(rs, model.Email)
	if err != nil {
		return 0, err
	}

	if exists == 1 {
		return 0, errors.New("User already exist")
	}

	s := New()
	model.Password = app.CalculatePassHash(model.Password, s.Salt)
	model.Salt = s.Salt
	verificationCode, _ := uuid.NewV4()
	model.VerificationCode = verificationCode.String()

	if err := u.dao.Register(rs, model); err != nil {
		return 0, err
	}

	//submit user role
	ur := models.AdminUserRoles{}
	ur.UserID = model.UserID
	ur.URID = 0
	if err := ur.ValidateRoles(); err != nil {
		return 0, err
	}

	if err := u.dao.SubmitUserRole(rs, &ur); err != nil {
		return 0, err
	}

	// Send Verification email
	// if err := u.sendConfirmationEmail(rs, model.UserID, model.Email, model.VerificationCode); err != nil {
	//	return 0, err
	//}

	return model.UserID, nil
}

// SubmitUserRole creates a new user role.
func (u *UserService) SubmitUserRole(rs app.RequestScope, model *models.AdminUserRoles) (*models.AdminUserRoles, error) {
	model.URID = app.GenerateNewID()

	if err := model.ValidateRoles(); err != nil {
		return nil, err
	}

	if err := u.dao.SubmitUserRole(rs, model); err != nil {
		return nil, err
	}

	return model, nil
}

// Delete deletes the user with the specified ID.
func (u *UserService) Delete(rs app.RequestScope, id int32) error {
	err := u.dao.Delete(rs, id)
	return err
}

// // sendConfirmationEmail Sends an email to the new registered user
// func (u *UserService) sendConfirmationEmail(rs app.RequestScope, userID uint64, recipientAddress string, verificationID string) error {

// 	//Construct verification url
// 	verificationURL := app.Config.VerificationLink + "?identity=" + strconv.Itoa(int(userID)) + "&confirm_verification=" + verificationID

// 	// Create a template using cds_confirmation_email.html
// 	emailData := models.ConfirmationEmailDetails{
// 		Title:            "Action Required: Please verify your email address.",
// 		InitialMessage:   "Thank you for creating a gPa Account. Please verify your email address to complete the registration process. Click the button below (it only takes a few seconds).",
// 		ButtonMessage:    "Verify your email address",
// 		FinalMessage:     "Verifying your email address ensures that you can securely retrieve your account information if your password is lost or stolen. You must verify your email address before you can use it on gPa services.",
// 		FinalMessage1:    "Thanks.",
// 		VerificationLink: verificationURL,
// 		UserID:           userID,
// 		Email:            recipientAddress,
// 		VerificationCode: verificationID,
// 	}

// 	if err := emailData.VerifyConfirmationEmail(); err != nil {
// 		return err
// 	}

// 	if verificationID == "" {
// 		verificationCode, _ := uuid.NewV4()
// 		emailData.VerificationCode = verificationCode.String()
// 		if err := u.dao.CreateNewEmailVerification(rs, &emailData); err != nil {
// 			return err
// 		}
// 	}

// 	absPath, _ := filepath.Abs("./views/cds_confirmation_email.html")

// 	tmpl, err := template.New("cds_confirmation_email.html").ParseFiles(absPath)
// 	if err != nil {
// 		return err
// 	}

// 	// Stores the parsed template
// 	var buff bytes.Buffer

// 	// Send the parsed template to buff
// 	err = tmpl.Execute(&buff, emailData)
// 	if err != nil {
// 		return err
// 	}

// 	body := buff.String()

// 	// Create the mail to send
// 	newMail := models.CreateMail(app.Config.DefaultMailSender, emailData.Email, emailData.Title, body)

// 	defer func() {
// 		if p := recover(); p != nil {
// 			err = p.(error)
// 		}
// 	}()

// 	// Send the mail just created
// 	newMail.Send()

// 	return nil
// }
