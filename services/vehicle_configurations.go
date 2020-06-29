package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
)

// vehicleDAO specifies the interface of the vehicle DAO needed by VehicleService.
type vehicleDAO interface {
	GetVehicleByStrID(rs app.RequestScope, strid string) (*models.VehicleConfigDetails, error)
	GetConfigurationDetails(rs app.RequestScope, vehicleid, deviceid int) (*models.VehicleConfigDetails, error)
	CountTripRecords(rs app.RequestScope, deviceid string) (int, error)
	GetTripDataByDeviceIDBtwDates(deviceid string, offset, limit int, from, to int64) ([]models.DeviceData, error)
	GetVehicleName(rs app.RequestScope, deviceid int) models.VDetails
	// Create saves a new vehicle in the storage.
	CreateVehicle(rs app.RequestScope, vehicle *models.VehicleDetails) (uint32, error)
	CreateVehicleOwner(rs app.RequestScope, vo *models.VehicleOwner) (uint32, error)
	CreateFitter(rs app.RequestScope, fd *models.FitterDetails) error
	CreateConfiguration(rs app.RequestScope, vehicle *models.Vehicle, ownerid uint32, fitterid uint32, vehicleid uint32, vehstringid string) error
	UpdateConfigurationStatus(rs app.RequestScope, configid uint32, status int8) error
	CountOverspeed(rs app.RequestScope, deviceid string) (int, error)
	CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error)
	GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.DeviceData, error)
	GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.DeviceData, error)
	DeleteOverspeedsByDeviceID(rs app.RequestScope, id uint32) (int, error)
	SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int, qtype string) ([]models.SearchDetails, error)
	CountSearches(rs app.RequestScope, searchterm, qtype string) (int, error)
	UpdatDeviceConfigurationStatus(rs app.RequestScope, deviceid int64, vehicleid uint32) error
	GetTripDataByDeviceID(deviceid string, offset, limit int, orderby string) ([]models.DeviceData, error)
	CountTripDataByDeviceID(deviceid string) (int, error)
	CountTripRecordsBtwDates(deviceid string, from int64, to int64) (int, error)
	CountAllViolations() (int, error)
	ListAllViolations(rs app.RequestScope, offset, limit int) ([]models.CurrentViolations, error)
	XMLListAllViolations(rs app.RequestScope, offset, limit int) ([]models.XMLResults, error)
	// CreateDevice saves a new device in the storage.
	CreateDevice(rs app.RequestScope, device *models.Devices) error
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

// GetConfigurationDetails ...
func (s *VehicleService) GetConfigurationDetails(rs app.RequestScope, vehicleid, deviceid int) (*models.VehicleConfigDetails, error) {
	return s.dao.GetConfigurationDetails(rs, vehicleid, deviceid)
}

// CountTripDataByDeviceID ...
func (s *VehicleService) CountTripDataByDeviceID(deviceid string) (int, error) {
	return s.dao.CountTripDataByDeviceID(deviceid)
}

// GetTripDataByDeviceID ...
func (s *VehicleService) GetTripDataByDeviceID(deviceid string, offset, limit int, orderby string) ([]models.DeviceData, error) {

	var deviceData []models.DeviceData
	var data []string
	var err error

	if orderby == "asc" {
		data, err = app.ZRevRange("data:"+deviceid, int64(offset), int64(limit+offset-1))
	} else {
		data, err = app.ZRevRange("data:"+deviceid, int64(offset), int64(limit+offset-1))
	}
	if err != nil {
		fmt.Println("Getting Keys Failed : " + err.Error())
	}

	if len(data) < limit {
		deviceData, err := s.dao.GetTripDataByDeviceID(deviceid, offset, limit, orderby)
		for _, rec := range deviceData {
			go app.ZAdd("data:"+deviceid, rec.DateTimeStamp, rec)
		}
		return deviceData, err
	}

	for i := 0; i < len(data); i++ {

		if data[i] != "0" {
			var deserializedValue models.DeviceData
			json.Unmarshal([]byte(data[i]), &deserializedValue)
			deviceData = append(deviceData, deserializedValue)
		}

	}

	return deviceData, err
}

// CountTripRecordsBtwDates ...
func (s *VehicleService) CountTripRecordsBtwDates(deviceid string, from, to int64) (int, error) {
	return s.dao.CountTripRecordsBtwDates(deviceid, from, to)
}

// GetTripDataByDeviceIDBtwDates ...
func (s *VehicleService) GetTripDataByDeviceIDBtwDates(deviceid string, offset, limit int, from, to int64) ([]models.DeviceData, error) {
	return s.dao.GetTripDataByDeviceIDBtwDates(deviceid, offset, limit, from, to)
}

// CountRedisTripRecords ...
func (s *VehicleService) CountRedisTripRecords(deviceid string) int {
	ListLength := app.ZCount("data:"+deviceid, "-inf", "+inf")
	return int(ListLength)
}

// GetOverspeedByDeviceID ...
func (s *VehicleService) GetOverspeedByDeviceID(rs app.RequestScope, deviceid string, offset, limit int) ([]models.DeviceData, error) {
	return s.dao.GetOverspeedByDeviceID(rs, deviceid, offset, limit)
}

// DeleteOverspeedsByDeviceID ...
func (s *VehicleService) DeleteOverspeedsByDeviceID(rs app.RequestScope, id uint32) (int, error) {
	return s.dao.DeleteOverspeedsByDeviceID(rs, id)
}

// GetViolationsByDeviceID ...
func (s *VehicleService) GetViolationsByDeviceID(rs app.RequestScope, deviceid string, reason string, offset, limit int) ([]models.DeviceData, error) {
	return s.dao.GetViolationsByDeviceID(rs, deviceid, reason, offset, limit)
}

// SearchVehicles ...
func (s *VehicleService) SearchVehicles(rs app.RequestScope, searchterm string, offset, limit int, qtype string) ([]models.SearchDetails, error) {
	return s.dao.SearchVehicles(rs, searchterm, offset, limit, qtype)
}

// GetCurrentViolations single violation as they happen...
func (s *VehicleService) GetCurrentViolations(rs app.RequestScope) (models.DeviceData, error) {
	var d models.DeviceData
	res, err := s.dao.ListAllViolations(rs, 0, 1)
	if err != nil {
		return d, err
	}
	fmt.Printf("Name is %v", res[0].Data.Name)

	// if value.SystemCode == "MCPG" {
	// 	var (
	// 		message          string
	// 		messageid        int
	// 		violationMessage = make(chan models.MessageDetails)
	// 	)

	// 	go app.SendViolationSMSMessages(violationMessage)

	// 	// fmt.Println("device_id", value.DeviceID)
	// 	if value.Offline {
	// 		message = value.Name + " offline at " + value.DateTime.Format(time.RFC3339)
	// 		messageid = 4
	// 	} else if value.Disconnect {
	// 		message = value.Name + " power disconnectd at " + value.DateTime.Format(time.RFC3339)
	// 		messageid = 3
	// 	} else if value.Failsafe {
	// 		message = value.Name + " signal disconnectd at " + value.DateTime.Format(time.RFC3339)
	// 		messageid = 2
	// 	} else if value.GroundSpeed > 80 {
	// 		message = value.Name + " was overspeeding at " + value.DateTime.Format(time.RFC3339)
	// 		messageid = 1
	// 	}

	// 	fmt.Println(messageid, message)
	// 	// violationMessage <- models.MessageDetails{messageid, message}
	// 	vd := s.dao.GetVehicleName(rs, int(value.DeviceID))
	// 	value.Name = vd.Name
	// 	deviceData = append(deviceData, value)
	// }

	return res[0].Data, err
}

// CountAllViolations ...
func (s *VehicleService) CountAllViolations() (int, error) {
	// count := app.ZCount("violations", "-inf", "+inf")
	// return int(count)
	return s.dao.CountAllViolations()
}

// ListAllViolations ...
func (s *VehicleService) ListAllViolations(rs app.RequestScope, offset, limit int) ([]models.CurrentViolations, error) {
	return s.dao.ListAllViolations(rs, offset, limit)
}

// XMLListAllViolations ...
func (s *VehicleService) XMLListAllViolations(rs app.RequestScope, offset, limit int) ([]models.XMLResults, error) {
	return s.dao.XMLListAllViolations(rs, offset, limit)
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
func (s *VehicleService) CountSearches(rs app.RequestScope, searchterm, qtype string) (int, error) {
	return s.dao.CountSearches(rs, searchterm, qtype)
}

// CountViolations Count returns the number of Violation records.
func (s *VehicleService) CountViolations(rs app.RequestScope, deviceid string, reason string) (int, error) {
	return s.dao.CountViolations(rs, deviceid, reason)
}

// Create creates a new vehicle configuration.
func (s *VehicleService) Create(rs app.RequestScope, model *models.Vehicle) (int, error) {
	// if err := model.Validate(); err != nil {
	//	return nil, err
	// }
	userid := model.UserID

	// Add Device Details
	did, _ := strconv.ParseInt(model.GovernorDetails.DeviceID, 10, 64)
	dm := models.NewDevice(did, model.DeviceDetails.DeviceType, model.DeviceDetails.SerialNO, model.SimNO, model.MotherboardNO, model.Technician)
	// if err := s.dao.CreateDevice(rs, dm); err != nil {
	// 	return 0, err
	// }
	s.dao.CreateDevice(rs, dm)

	// Add vehicle owner
	vm := NewOwner(model.DeviceDetails, model.OwnerID, userid)
	ownerid, err := s.dao.CreateVehicleOwner(rs, vm)
	if err != nil {
		return 0, err
	}

	// Add Vehicle
	vd := NewVehicle(model.DeviceDetails, model.VehicleID, userid)
	vehid, err := s.dao.CreateVehicle(rs, vd)
	if err != nil {
		return 0, err
	}
	model.VehicleID = vehid

	// Update Device Configuration status (is configured)
	if err := s.dao.UpdatDeviceConfigurationStatus(rs, did, model.VehicleID); err != nil {
		return 0, err
	}

	if err = s.dao.CreateConfiguration(rs, model, ownerid, model.FitterID, model.VehicleID, vd.VehicleStringID); err != nil {
		return 0, err
	}

	// Add vehicle to tracking server
	tsv := NewTrackingServerVehicle(model)
	_, err = AddDevicesTrackingServer(rs, tsv, "en", model.UserHash)
	if err != nil {
		return 0, err
	}

	return 0, nil
}
