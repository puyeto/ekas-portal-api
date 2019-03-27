package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// trackingServerServiceDAO specifies the interface of the trackingServerService DAO needed by TrackingServerServiceService.
type trackingServerServiceDAO interface {
}

// TrackingServerServiceService provides services related with trackingServerServices.
type TrackingServerServiceService struct {
	dao trackingServerServiceDAO
}

// NewTrackingServerServiceService creates a new TrackingServerServiceService with the given trackingServerService DAO.
func NewTrackingServerServiceService(dao trackingServerServiceDAO) *TrackingServerServiceService {
	return &TrackingServerServiceService{dao}
}

var serviceTypes = models.TrackingServiceTypes{
	Days:        "Days",
	Odometer:    "Odometer",
	EngineHours: "engine_hours",
}

// TrackingServerGetServices ...
func (s *TrackingServerServiceService) TrackingServerGetServices(rs app.RequestScope, model *models.UserData) (interface{}, error) {
	p := url.Values{
		"user_api_hash": {model.UserHash},
		"lang":          {model.Lang},
		"device_id":     {string(model.DeviceID)},
	}

	URL := app.Config.TrackingServerURL + "add_service_data?" + p.Encode()
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

	return data, nil
}

// TrackingServerAddServices ...
func (s *TrackingServerServiceService) TrackingServerAddServices(rs app.RequestScope, model *models.AddTrackingServiceDetails, lang string, userhash string, deviceid string) (interface{}, error) {

	p := url.Values{
		"user_api_hash":          {userhash},
		"lang":                   {lang},
		"device_id":              {deviceid},
		"name":                   {model.Name},
		"expiration_by":          {model.ExpirationBy},
		"interval":               {model.Interval},
		"last_service":           {model.LastService},
		"trigger_event_left":     {model.TriggerEventLeft},
		"renew_after_expiration": {string(model.RenewAfterExpiration)},
		"email":                  {model.Email},
		"mobile_phone":           {model.MobilePhone},
	}

	URL := app.Config.TrackingServerURL + "add_service?" + p.Encode()

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

	return data, nil
}
