package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	//"time"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// vehicleDAO specifies the interface of the vehicle DAO needed by VehicleService.
type vehicleDAO interface {
	GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error)
	GetTripDataByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.TripData, error)
	CountTripRecords(rs app.RequestScope, deviceid string) (int, error)
	FetchAllTripsBetweenDates(rs app.RequestScope, deviceid string, offset, limit int, from string, to string) ([]models.TripData, error)
	ListRecentViolations(rs app.RequestScope, offset, limit int) ([]models.CurrentViolations, error)
	// Create saves a new vehicle in the storage.
	CreateVehicle(rs app.RequestScope, vehicle *models.VehicleDetails) error
	CreateVehicleOwner(rs app.RequestScope, vo *models.VehicleOwner) error
	CreateFitter(rs app.RequestScope, fd *models.FitterDetails) error
	CreateConfiguration(rs app.RequestScope, vehicle *models.Vehicle, ownerid uint32, fitterid uint32, vehicleid uint32) error
	UpdateConfigurationStatus(rs app.RequestScope, configid uint32, status int8) error
	CountOverspeed(rs app.RequestScope, deviceid string) (int, error)
	CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error)
	GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.TripData, error)
	GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.TripData, error)
	SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int) ([]models.SearchDetails, error)
	CountSearches(rs app.RequestScope, searchterm string) (int, error)
}

// VehicleService provides services related with vehicles.
type VehicleService struct {
	dao vehicleDAO
}

// NewVehicleService creates a new VehicleService with the given vehicle DAO.
func NewVehicleService(dao vehicleDAO) *VehicleService {
	return &VehicleService{dao}
}

// GetVehicleByStrID ...
func (s *VehicleService) GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error) {
	return s.dao.GetVehicleByStrID(rs, strid)
}

// GetTripDataByDeviceID ...
func (s *VehicleService) GetTripDataByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.DeviceData, error) {
	// return s.dao.GetTripDataByDeviceID(rs, deviceid, offset, limit)
	// define slice of Identification
	var deviceData []models.DeviceData

	keysList, err := app.ZRevRange("data:"+deviceid, int64(offset), int64(limit))
	if err != nil {
		fmt.Println("Getting Keys Failed : " + err.Error())
	}

	for i := 0; i < len(keysList); i++ {

		if keysList[i] != "0" {
			var deserializedValue models.DeviceData
			json.Unmarshal([]byte(keysList[i]), &deserializedValue)
			deviceData = append(deviceData, deserializedValue)
		}

	}

	return deviceData, err
}

// FetchAllTripsBetweenDates ...
func (s *VehicleService) FetchAllTripsBetweenDates(rs app.RequestScope, deviceid string, offset, limit int, from, to int64) ([]models.DeviceData, error) {
	var deviceData []models.DeviceData

	min := strconv.FormatInt(from, 10)
	max := strconv.FormatInt(to, 10)

	keysList, err := app.ZRevRangeByScore("data:"+deviceid, min, max, int64(offset), int64(limit))
	if err != nil {
		fmt.Println("Getting Keys Failed : " + err.Error())
	}

	for i := 0; i < len(keysList); i++ {

		if keysList[i] != "0" {
			var deserializedValue models.DeviceData
			json.Unmarshal([]byte(keysList[i]), &deserializedValue)
			deviceData = append(deviceData, deserializedValue)
		}

	}

	return deviceData, err
}

// CountRedisTripRecords ...
func (s *VehicleService) CountRedisTripRecords(rs app.RequestScope, deviceid string) int {
	ListLength := app.ZCount("data:"+deviceid, "-inf", "+inf")
	return int(ListLength)
}

// CountRedisTripRecordsBtwDates ...
func (s *VehicleService) CountRedisTripRecordsBtwDates(rs app.RequestScope, deviceid string, from, to int64) int {
	min := strconv.FormatInt(from, 10)
	max := strconv.FormatInt(to, 10)
	count := app.ZCount("data:"+deviceid, min, max)
	return int(count)
}

// GetOverspeedByDeviceID ...
func (s *VehicleService) GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.TripData, error) {
	return s.dao.GetOverspeedByDeviceID(rs, deviceid, offset, limit)
}

// GetViolationsByDeviceID ...
func (s *VehicleService) GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.TripData, error) {
	return s.dao.GetViolationsByDeviceID(rs, deviceid, reason, offset, limit)
}

// SearchVehicles ...
func (s *VehicleService) SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int) ([]models.SearchDetails, error) {
	return s.dao.SearchVehicles(rs, searchterm, offset, limit)
}

// ListRecentViolations ...
func (s *VehicleService) ListRecentViolations(rs app.RequestScope, offset, limit int) ([]models.CurrentViolations, error) {
	return s.dao.ListRecentViolations(rs, offset, limit)
}

// GetCurrentViolations single violation as they happen...
func (s *VehicleService) GetCurrentViolations(rs app.RequestScope) ([]models.DeviceData, error) {
	// define slice of Identification
	var deviceData []models.DeviceData

	value, err := app.GetDeviceDataValue("currentviolations")
	if err != nil {
		return nil, err
	}
	if value.SystemCode == "MCPG" {
		// var message string
		// var message_id int
		// // fmt.Println("device_id", value.DeviceID)
		// if value.Offline {
		// 	message = value.Name + " offline at " + value.DateTime.Format(time.RFC3339)
		// 	message_id = 4
		// } else if value.Disconnect {
		// 	message = value.Name + " power disconnectd at " + value.DateTime.Format(time.RFC3339)
		// 	message_id = 3
		// } else if value.Failsafe {
		// 	message = value.Name + " signal disconnectd at " + value.DateTime.Format(time.RFC3339)
		// 	message_id = 2
		// } else if value.GroundSpeed > 80 {
		// 	message = value.Name + " was overspeeding at " + value.DateTime.Format(time.RFC3339)
		// 	message_id = 1
		// }

		// app.Message <- models.MessageDetails{message_id, message}

		// go app.SendSMSMessages("+254723436438", message)
		deviceData = append(deviceData, value)
	}

	return deviceData, err
}

// ListAllViolations ...
func (s *VehicleService) ListAllViolations(rs app.RequestScope, offset, limit int) ([]models.DeviceData, error) {

	var deviceData []models.DeviceData

	keysList, err := app.ZRevRange("violations", int64(offset), int64(limit))
	if err != nil {
		fmt.Println("Getting Keys Failed : " + err.Error())
	}

	for i := 0; i < len(keysList); i++ {

		if keysList[i] != "0" {
			var deserializedValue models.DeviceData
			json.Unmarshal([]byte(keysList[i]), &deserializedValue)
			deviceData = append(deviceData, deserializedValue)
		}

	}

	return deviceData, err
}

// GetOfflineViolations ...
func (s *VehicleService) GetOfflineViolations(rs app.RequestScope, deviceid string) ([]models.DeviceData, error) {

	var deviceData []models.DeviceData

	keysList, err := app.ZRevRange("offline:"+deviceid, 0, 100)
	if err != nil {
		fmt.Println("Getting Keys Failed : " + err.Error())
	}

	for i := 0; i < len(keysList); i++ {

		if keysList[i] != "0" {
			var deserializedValue models.DeviceData
			json.Unmarshal([]byte(keysList[i]), &deserializedValue)
			deviceData = append(deviceData, deserializedValue)
		}

	}

	return deviceData, err
}

// CountAllViolations ...
func (s *VehicleService) CountAllViolations(rs app.RequestScope) int {
	count := app.ZCount("violation", "-inf", "+inf")
	return int(count)
}

// GetUnavailableDevices ...
func (s *VehicleService) GetUnavailableDevices(rs app.RequestScope) ([]models.DeviceData, error) {

	// define slice of Identification
	var deviceData []models.DeviceData

	keysList, err := app.ListKeys("lastseen:*")
	if err != nil {
		fmt.Println("Getting Keys Failed : " + err.Error())
	}

	for i := 0; i < len(keysList); i++ {
		fmt.Println("Getting " + keysList[i])
		value, err := app.GetLastSeenValue(keysList[i])
		if err != nil {
			return nil, err
		}
		if value.SystemCode == "MCPG" {
			if callTime(value) >= 5 {
				fmt.Println("device_id", value.DeviceID)
				deviceData = append(deviceData, value)
			}
		}
	}

	return deviceData, err
}

func callTime(m models.DeviceData) int {
	nowd := time.Now()
	now := dateF(nowd.Year(), nowd.Month(), nowd.Day(), nowd.Hour(), nowd.Minute(), nowd.Second())
	pastDate := dateF(m.UTCTimeYear, time.Month(m.UTCTimeMonth), m.UTCTimeDay, m.UTCTimeHours, m.UTCTimeMinutes, m.UTCTimeSeconds)
	diff := now.Sub(pastDate)

	mins := int(diff.Minutes())
	fmt.Println("mins = ", mins)
	return mins
}

func dateF(year int, month time.Month, day int, hr, min, sec int) time.Time {
	return time.Date(year, month, day, hr, min, sec, 0, time.UTC)
}

// CountTripRecords Count returns the number of trip records.
func (s *VehicleService) CountTripRecords(rs app.RequestScope, deviceid string) (int, error) {
	return s.dao.CountTripRecords(rs, deviceid)
}

// CountOverspeed Count returns the number of overspeed records.
func (s *VehicleService) CountOverspeed(rs app.RequestScope, deviceid string) (int, error) {
	return s.dao.CountOverspeed(rs, deviceid)
}

// CountSearches Count returns the number of search records.
func (s *VehicleService) CountSearches(rs app.RequestScope, searchterm string) (int, error) {
	return s.dao.CountSearches(rs, searchterm)
}

// CountViolations Count returns the number of Violation records.
func (s *VehicleService) CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error) {
	return s.dao.CountViolations(rs, deviceid, reason)
}

// Create creates a new vehicle.
func (s *VehicleService) Create(rs app.RequestScope, model *models.Vehicle) (int, error) {
	// if err := model.Validate(); err != nil {
	//	return nil, err
	// }

	// Add vehicle owner
	var ownerid = app.GenerateNewID()
	if model.OwnerID > 0 {
		ownerid = model.OwnerID
	}
	vm := NewOwner(model.DeviceDetails, ownerid)
	if err := s.dao.CreateVehicleOwner(rs, vm); err != nil {
		return 0, err
	}

	// Add Fitter Center / Fitter
	var fid = app.GenerateNewID()
	if model.FitterID > 0 {
		fid = model.FitterID
	}
	fd := NewFitter(model.DeviceDetails, fid)
	if err := s.dao.CreateFitter(rs, fd); err != nil {
		return 0, err
	}

	// Add Vehicle
	fmt.Println(model.VehicleID)
	var vid = app.GenerateNewID()
	if model.VehicleID > 0 {
		vid = model.VehicleID
	}
	fmt.Println(vid)
	vd := NewVehicle(model.DeviceDetails, vid)
	if err := s.dao.CreateVehicle(rs, vd); err != nil {
		return 0, err
	}

	// Add Configuartion Details
	if model.ConfigID > 0 {
		// update configuration status
		if err := s.dao.UpdateConfigurationStatus(rs, model.ConfigID, 0); err != nil {
			return 0, err
		}
	}
	if err := s.dao.CreateConfiguration(rs, model, vm.OwnerID, fd.FitterID, vd.VehicleID); err != nil {
		return 0, err
	}

	// Add vehicle to tracking server
	tsv := NewTrackingServerVehicle(model)
	_, err := AddDevicesTrackingServer(rs, tsv, "en", model.UserHash)
	if err != nil {
		return 0, err
	}

	return 0, nil
}
