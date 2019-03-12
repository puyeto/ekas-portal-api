package daos

import (
	"encoding/json"

	"github.com/ekas-portal-api/app"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// TrackingServerDAO persists trackingServer data in database
type TrackingServerDAO struct{}

// NewTrackingServerDAO creates a new TrackingServerDAO
func NewTrackingServerDAO() *TrackingServerDAO {
	return &TrackingServerDAO{}
}

// GetTrackingServerUserLoginIDByEmail ...
func (dao *TrackingServerDAO) GetTrackingServerUserLoginIDByEmail(rs app.RequestScope, email string) (uint64, error) {
	var uid uint64
	q := rs.Tx().NewQuery("SELECT auth_user_id FROM auth_users WHERE auth_user_email='" + email + "' LIMIT 1")
	err := q.Row(&uid)
	return uid, err
}

// SaveTrackingServerLoginDetails saves a new user record in the database.
func (dao *TrackingServerDAO) SaveTrackingServerLoginDetails(rs app.RequestScope, id uint64, email string, hash string, status int8, data interface{}) error {
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

// TrackingServerUserEmailExists check if a tracker server auth user exists
func (dao *TrackingServerDAO) TrackingServerUserEmailExists(rs app.RequestScope, email string) (int, error) {
	var exists int
	q := rs.Tx().NewQuery("SELECT EXISTS(SELECT 1 FROM auth_users WHERE auth_user_email='" + email + "' LIMIT 1) AS exist")
	err := q.Row(&exists)
	return exists, err
}
