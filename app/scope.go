package app

import (
	"net/http"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/sirupsen/logrus"
)

// RequestScope contains the application-specific information that are carried around in a request.
type RequestScope interface {
	Logger
	// UserID returns the ID of the user for the current request
	UserID() string
	// SetUserID sets the ID of the currently authenticated user
	SetUserID(id string)
	CompanyID() string
	SetCompanyID(id string)
	RoleID() string
	SetRoleID(id string)
	// RequestID returns the ID of the current request
	RequestID() string
	// Tx returns the currently active database transaction that can be used for DB query purpose
	Tx() *dbx.Tx
	// SetTx sets the database transaction
	SetTx(tx *dbx.Tx)
	// Rollback returns a value indicating whether the current database transaction should be rolled back
	Rollback() bool
	// SetRollback sets a value indicating whether the current database transaction should be rolled back
	SetRollback(bool)
	// Now returns the timestamp representing the time when the request is being processed
	Now() time.Time
}

type requestScope struct {
	Logger              // the logger tagged with the current request information
	now       time.Time // the time when the request is being processed
	requestID string    // an ID identifying one or multiple correlated HTTP requests
	userID    string    // an ID identifying the current user
	roleID    string    // an ID identifying the current user
	companyID string
	rollback  bool    // whether to roll back the current transaction
	tx        *dbx.Tx // the currently active transaction
}

func (rs *requestScope) UserID() string {
	return rs.userID
}

func (rs *requestScope) CompanyID() string {
	return rs.companyID
}

func (rs *requestScope) RoleID() string {
	return rs.roleID
}

func (rs *requestScope) SetUserID(id string) {
	rs.Logger.SetField("UserID", id)
	rs.userID = id
}

func (rs *requestScope) SetCompanyID(id string) {
	rs.Logger.SetField("CompanyID", id)
	rs.companyID = id
}

func (rs *requestScope) SetRoleID(id string) {
	rs.Logger.SetField("RoleID", id)
	rs.roleID = id
}

func (rs *requestScope) RequestID() string {
	return rs.requestID
}

func (rs *requestScope) Tx() *dbx.Tx {
	return rs.tx
}

func (rs *requestScope) SetTx(tx *dbx.Tx) {
	rs.tx = tx
}

func (rs *requestScope) Rollback() bool {
	return rs.rollback
}

func (rs *requestScope) SetRollback(v bool) {
	rs.rollback = v
}

func (rs *requestScope) Now() time.Time {
	return rs.now
}

// newRequestScope creates a new RequestScope with the current request information.
func newRequestScope(now time.Time, logger *logrus.Logger, request *http.Request) RequestScope {
	l := NewLogger(logger, logrus.Fields{})
	requestID := request.Header.Get("X-Request-Id")
	if requestID != "" {
		l.SetField("RequestID", requestID)
	}
	return &requestScope{
		Logger:    l,
		now:       now,
		requestID: requestID,
	}
}
