package models

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	gomail "gopkg.in/gomail.v2"
)

var (
	// ErrMissingField missing error message
	ErrMissingField = "Error missing %v"
)

// Credential user login credentials
type Credential struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

//UserDetails user structure
type AdminUserDetails struct {
	UserID           int32 `json:"user_id" db:"pk,user_id"`
	FirstName        string `json:"user_first_name" db:"first_name"`
	LastName         string `json:"user_last_name" db:"last_name"`
	Email            string `json:"user_email" db:"email"`
	Username         string `json:"username,omitempty" db:"username"`
	Password         string `json:"user_password,omitempty" db:"password"`
	DOB              string `json:"user_dob,omitempty" db:"dob"`
	MobileNumber     string `json:"user_mobile_number,omitempty" db:"mobile_number"`
	Salt             string `json:"salt,omitempty" db:"salt"`
	VerificationCode string `json:"Verification_code,omitempty" db:"Verification_code"`
	Token            string `json:"token,omitempty"`
	IsVerified       int8   `json:"is_verified,omitempty"`
}

// ListUserDetails list users structure
type ListUserDetails struct {
	UserDetails AdminUserDetails `json:"user_details,omitempty" db:"user_details"`
	Roles       []AdminUserRoles `json:"user_roles,omitempty" db:"roles"`
}

// CreateUser -
type CreateUser struct {
	UserDetails *AdminUserDetails `json:"user_details,omitempty" db:"user_details"`
	Roles       *AdminUserRoles   `json:"user_roles,omitempty" db:"user_roles"`
}

// ValidateCredential validates the login fields.
func (c Credential) ValidateCredential() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(6, 120)),
	)
}

//ValidateNewUser validate create user
func (u AdminUserDetails) ValidateNewUser() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(6, 120)),
	)
}

// NewSalt generate new salt
func (u *AdminUserDetails) NewSalt() {
	h := sha1.New()
	io.WriteString(h, strconv.Itoa(int(time.Now().UnixNano())))
	u.Salt = fmt.Sprintf("%x", h.Sum(nil))
}

// AdminUserRoles user role structure
type AdminUserRoles struct {
	URID     uint32 `json:"ur_id" db:"pk,ur_id"`
	UserID   int32  `json:"user" db:"user_id"`
	RoleID   int32  `json:"role" db:"role_id"`
	RoleName string `json:"role_name" db:"role_name"`
}

// ValidateRoles validate create user role
func (u AdminUserRoles) ValidateRoles() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.UserID, validation.Required),
	)
}

// ConfirmationEmailDetails ---
type ConfirmationEmailDetails struct {
	UserID           uint64 `valid:"required"`
	Email            string `valid:"required,email"`
	VerificationCode string `valid:"required"`
	Title            string
	InitialMessage   string
	ButtonMessage    string
	FinalMessage     string
	FinalMessage1    string
	VerificationLink string
}

// VerifyConfirmationEmail ...
func (con ConfirmationEmailDetails) VerifyConfirmationEmail() error {
	return validation.ValidateStruct(&con,
		validation.Field(&con.Email, validation.Required, is.Email),
		validation.Field(&con.UserID, validation.Required),
	)
}

// UserLoginSessions ...
type UserLoginSessions struct {
	SessionID uint32 `json:"session_id" db:"pk,session_id"`
	UserID    int32  `json:"user_id" db:"user_id"`
	UserAgent string `json:"user_agent" db:"user_agent"`
	IP        string `json:"ip" db:"ip"`
	Token     string `json:"token" db:"token"`
}

// MailDetails ---
type MailDetails struct {
	From    string
	To      string
	Subject string
	Body    string
}

// CreateMail creates a new mail
func CreateMail(from string, to string, subject string, body string) *MailDetails {
	return &MailDetails{
		From:    from,
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

//GetSender - Get sender email address
func (m *MailDetails) GetSender() string {
	return m.From
}

//GetRecipient - Get recipient email address
func (m *MailDetails) GetRecipient() string {
	return m.To
}

//GetSubject - Get the email's subject
func (m *MailDetails) GetSubject() string {
	return m.Subject
}

//GetBody - Get the email's body
func (m *MailDetails) GetBody() string {
	return m.Body
}

//Send - sends out the mail
func (m *MailDetails) Send() {
	s := gomail.NewMessage()
	s.SetHeader("From", m.GetSender())
	s.SetHeader("To", m.GetRecipient())
	s.SetHeader("Subject", m.GetSubject())
	s.SetBody("text/html", m.GetBody())

	MailChannel <- s

}

//MailChannel - variable
var MailChannel chan *gomail.Message

// MailDaemon mail daemon listening for mails to send
// func MailDaemon(ch chan *gomail.Message) {
// 	go func() {
// 		d := gomail.NewDialer(app.Config.MailHost, app.Config.MailPort, app.Config.MailUsername, app.Config.MailPassword)
// 		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

// 		var s gomail.SendCloser
// 		var err error
// 		open := false
// 		for {
// 			select {
// 			case m, ok := <-ch:
// 				if !ok {
// 					return
// 				}
// 				if !open {
// 					if s, err = d.Dial(); err != nil {
// 						// panic(err)
// 						log.Print(err)
// 					}
// 					open = true
// 				}
// 				if err := gomail.Send(s, m); err != nil {
// 					log.Print(err)
// 				}
// 			// Close the connection to the SMTP server if no email was sent in
// 			// the last 30 seconds.
// 			case <-time.After(30 * time.Second):
// 				if open {
// 					if err := s.Close(); err != nil {
// 						panic(err)
// 					}
// 					open = false
// 				}
// 			}
// 		}
// 	}()
// }

// GetRemoteIP ---
func GetRemoteIP(r *http.Request) string {
	fwdIP := r.Header.Get("X-Forwarded-For")
	fwdSplit := strings.Split(fwdIP, ",")
	if fwdIP != "" {
		// pick the leftmost x-forwarded-for addr
		return fwdSplit[0]
	}

	// this literally can't fail on r.RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
