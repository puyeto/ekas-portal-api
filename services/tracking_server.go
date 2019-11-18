package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// trackingServerDAO specifies the interface of the trackingServer DAO needed by TrackingServerService.
type trackingServerDAO interface {
	// Login to tracking server.
	SaveTrackingServerLoginDetails(rs app.RequestScope, id uint32, email string, hash string, status int8, data interface{}) error
	TrackingServerUserEmailExists(rs app.RequestScope, email string) (int, error)
	GetTrackingServerUserLoginIDByEmail(rs app.RequestScope, email string) (uint32, int, int, error)
	CreateLoginSession(rs app.RequestScope, ls *models.UserLoginSessions) error
	GetUserByEmail(rs app.RequestScope, email string) (models.AdminUserDetails, error)
	QueryVehicelsFromPortal(rs app.RequestScope, offset, limit int, uid int) ([]models.VehicleDetails, error)
	GetUserByUserHash(rs app.RequestScope, userhash string) (models.AdminUserDetails, error)
}

// TrackingServerService ---
type TrackingServerService struct {
	dao trackingServerDAO
}

// NewTrackingServerService creates a new TrackingServerService with the given artist DAO.
func NewTrackingServerService(dao trackingServerDAO) *TrackingServerService {
	return &TrackingServerService{dao}
}

// TrackingServerLogin ...
func (s *TrackingServerService) TrackingServerLogin(rs app.RequestScope, model *models.TrackingServerAuth) (m models.AdminUserDetails, err error) {
	if err := model.ValidateTrackingServerLogin(); err != nil {
		return m, err
	}

	return s.Login(rs, model.Email, model.Password)
}

// TrackingServerLogin2 login to the tracking server
// func (s *TrackingServerService) TrackingServerLogin2(rs app.RequestScope, model *models.TrackingServerAuth) (models.AdminUserDetails, error) {
// 	// if err := model.ValidateTrackingServerLogin(); err != nil {
// 	// 	return nil, err
// 	// }
// 	URL := app.Config.TrackingServerURL + "login/?email=" + model.Email + "&password=" + model.Password
// 	res, err := http.Get(URL)
// 	if err != nil {
// 		// return nil, err
// 		return s.Login(rs, model.Email, model.Password)
// 	}

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var data map[string]interface{}
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var id = app.GenerateNewID()
// 	var hash = data["user_api_hash"].(string)
// 	var status = int8(data["status"].(float64))
// 	data["user_email"] = model.Email
// 	data["user_id"] = id
// 	data["user_role"] = 10005

// 	exists, err := s.dao.TrackingServerUserEmailExists(rs, model.Email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if exists == 1 {
// 		uid, role, cid, err := s.dao.GetTrackingServerUserLoginIDByEmail(rs, model.Email)
// 		data["user_id"] = uid
// 		data["user_role"] = role
// 		data["company_id"] = cid
// 		if err != nil {
// 			return nil, err
// 		}
// 	} else {
// 		// Save Results to db
// 		err = s.dao.SaveTrackingServerLoginDetails(rs, id, model.Email, hash, status, data)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	return data, nil
// }

// Login a user  from portal
func (s *TrackingServerService) Login(rs app.RequestScope, email, password string) (models.AdminUserDetails, error) {

	res, err := s.dao.GetUserByEmail(rs, email)
	if err != nil {
		// return res, err
	}

	// reset(res)

	// token, err := auth.NewJWT(jwt.MapClaims{
	// 	"id":  strconv.Itoa(int(res.UserID)),
	// 	"exp": time.Now().Add(time.Hour * 72).Unix(),
	// }, app.Config.JWTSigningKey)
	// if err != nil {
	// 	return nil, errors.New(err.Error())
	// }

	s.storeLoginSession(rs, &res)

	return res, nil
}

// storeLoginSession ...
func (s *TrackingServerService) storeLoginSession(rs app.RequestScope, ud *models.AdminUserDetails) error {
	r := &http.Request{}
	loginSession := models.UserLoginSessions{
		SessionID: app.GenerateNewID(),
		UserID:    ud.UserID,
		UserAgent: r.UserAgent(),
		IP:        models.GetRemoteIP(r),
		Token:     ud.Token,
	}

	return s.dao.CreateLoginSession(rs, &loginSession)
}

// TrackingServerUserDevices - Get user devices from  the tracking server
func (s *TrackingServerService) TrackingServerUserDevices(rs app.RequestScope, model *models.UserData) ([]models.VehicleDetails, error) {
	res, err := s.dao.GetUserByUserHash(rs, model.UserHash)
	if err != nil {
		return nil, err
	}
	userid := res.UserID
	if res.Email == "ntsa@ekastech.com" {
		userid = 0
	}

	// get user id by UserHash
	vList, err := s.QueryVehicelsFromPortal(rs, 0, 100, int(userid))

	return vList, err
}

// func (s *TrackingServerService) TrackingServerUserDevices(rs app.RequestScope, model *models.UserData) (interface{}, error) {
// 	URL := app.Config.TrackingServerURL + "get_devices/?lang=" + model.Lang + "&user_api_hash=" + model.UserHash
// 	res, err := http.Get(URL)
// 	if err != nil {

// 		res, err := s.dao.GetUserByUserHash(rs, model.UserHash)
// 		if err != nil {
// 			return nil, err
// 		}
// 		userid := res.UserID
// 		if res.Email == "ntsa.ekastech.com" {
// 			userid = 0
// 		}
// 		// get user id by UserHash
// 		vList, err := s.QueryVehicelsFromPortal(rs, 0, 100, int(userid))

// 		return vList, err
// 	}

// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var data interface{}
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// fmt.Printf("Results: %v\n", data)

// 	return data, nil
// }

// QueryVehicelsFromPortal returns the vehicleRecords with the specified offset and limit.
func (s *TrackingServerService) QueryVehicelsFromPortal(rs app.RequestScope, offset, limit int, uid int) ([]models.VehicleDetails, error) {
	return s.dao.QueryVehicelsFromPortal(rs, offset, limit, uid)
}

// TrackingServerAddDevices - add user devices from  the tracking server
func (s *TrackingServerService) TrackingServerAddDevices(rs app.RequestScope, model *models.AddDeviceDetails, lang string, userhash string) (interface{}, error) {
	return AddDevicesTrackingServer(rs, model, lang, userhash)
}

// AddDevicesTrackingServer ...
func AddDevicesTrackingServer(rs app.RequestScope, model *models.AddDeviceDetails, lang string, userhash string) (interface{}, error) {
	p := url.Values{
		"user_api_hash":       {userhash},
		"lang":                {lang},
		"name":                {model.Name},
		"imei":                {model.Imei},
		"icon_id":             {model.IconID},
		"fuel_measurement_id": {model.FuelMeasurementID},
		"tail_length":         {model.TailLength},
		"min_fuel_thefts":     {model.MinFuelThefts},
		"min_moving_speed":    {model.MinFuelThefts},
		"min_fuel_fillings":   {model.MinFuelFillings},
		"plate_number":        {model.PlateNumber},
		"vin":                 {model.Vin},
		"device_model":        {model.DeviceModel},
		"registration_number": {model.RegistrationNumber},
		"object_owner":        {model.ObjectOwner},
		"device_icons_type":   {"icon"},
		"icon_stopped":        {"green"},
		"icon_offline":        {"orange"},
	}
	URL := app.Config.TrackingServerURL + "add_device?" + p.Encode()

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
	// fmt.Printf("Results: %v\n", data)

	return data, nil

}

// TrackingServerEditDevices - Edit user devices from the tracking server
func (s *TrackingServerService) TrackingServerEditDevices(rs app.RequestScope, model *models.AddDeviceDetails, lang string, userhash string) (interface{}, error) {

	p := url.Values{
		"user_api_hash":       {userhash},
		"lang":                {lang},
		"name":                {model.Name},
		"imei":                {model.Imei},
		"icon_id":             {model.IconID},
		"fuel_measurement_id": {model.FuelMeasurementID},
		"tail_length":         {model.TailLength},
		"min_fuel_thefts":     {model.MinFuelThefts},
		"min_moving_speed":    {model.MinFuelThefts},
		"min_fuel_fillings":   {model.MinFuelFillings},
		"plate_number":        {model.PlateNumber},
		"vin":                 {model.Vin},
		"device_model":        {model.DeviceModel},
		"registration_number": {model.RegistrationNumber},
		"object_owner":        {model.ObjectOwner},
		"device_icons_type":   {"icon"},
		"icon_stopped":        {"green"},
		"icon_offline":        {"orange"},
	}
	URL := app.Config.TrackingServerURL + "edit_device?" + p.Encode()

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
	// t.Printf("Results: %v\n", data)

	return data, nil
}
