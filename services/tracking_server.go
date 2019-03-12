package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// trackingServerDAO specifies the interface of the trackingServer DAO needed by TrackingServerService.
type trackingServerDAO interface {
	// Login to tracking server.
	SaveTrackingServerLoginDetails(rs app.RequestScope, id uint64, email string, hash string, status int8, data interface{}) error
	TrackingServerUserEmailExists(rs app.RequestScope, email string) (int, error)
	GetTrackingServerUserLoginIDByEmail(rs app.RequestScope, email string) (uint64, error)
}

// TrackingServerService ---
type TrackingServerService struct {
	dao trackingServerDAO
}

// NewTrackingServerService creates a new TrackingServerService with the given artist DAO.
func NewTrackingServerService(dao trackingServerDAO) *TrackingServerService {
	return &TrackingServerService{dao}
}

// TrackingServerLogin login to the tracking server
func (s *TrackingServerService) TrackingServerLogin(rs app.RequestScope, model *models.TrackingServerAuth) (interface{}, error) {
	if err := model.ValidateTrackingServerLogin(); err != nil {
		return nil, err
	}
	URL := app.Config.TrackingServerURL + "login/?email=" + model.Email + "&password=" + model.Password
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var id = app.GenerateNewID()
	var hash = data["user_api_hash"].(string)
	var status = int8(data["status"].(float64))
	data["user_email"] = model.Email
	data["user_id"] = id

	exists, err := s.dao.TrackingServerUserEmailExists(rs, model.Email)
	if err != nil {
		return nil, err
	}

	if exists == 1 {
		id, err := s.dao.GetTrackingServerUserLoginIDByEmail(rs, model.Email)
		data["user_id"] = id
		if err != nil {
			return nil, err
		}
	} else {
		// Save Results to db
		err = s.dao.SaveTrackingServerLoginDetails(rs, id, model.Email, hash, status, data)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// TrackingServerUserDevices - Get user devices from  the tracking server
func (s *TrackingServerService) TrackingServerUserDevices(rs app.RequestScope, model *models.UserData) (interface{}, error) {
	URL := app.Config.TrackingServerURL + "get_devices/?lang=" + model.Lang + "&user_api_hash=" + model.UserHash
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Results: %v\n", data)

	return data, nil
}
