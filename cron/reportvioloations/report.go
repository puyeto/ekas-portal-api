package reportvioloations

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Status ...
type Status struct {
	// filtered
}

// Run LastSeen.Run() will get triggered automatically.
func (c Status) Run() {
	// select all deviceids
	getAllViolations(0, 50000)
}

func sendSMS(vid, did int32, message, tonumber, violtype string) {
	check, _ := app.CheckMessages(tonumber, "Violations")
	duration := time.Now().Sub(check.DateTime)
	// fmt.Printf("difference %d days", int(duration.Hours()/24) )
	if int(duration.Hours()/24) < 4 {
		return
	}

	app.MessageChan <- app.MessageDetails{
		Message:  message,
		ToNumber: tonumber,
		Type:     violtype,
	}

	// save sms
	go saveSMS(vid, did, message, tonumber, violtype)

}

func saveSMS(vid, did int32, message, tonumber, violtype string) {
	details := models.SaveMessageDetails{
		MessageID:   0,
		Message:     message,
		MessageType: violtype,
		VehicleID:   vid,
		DeviceID:    did,
		DateTime:    time.Now(),
		Status:      "Sent",
		From:        "EKASTECH",
		To:          tonumber,
	}
	app.SaveSentMessages(details)
}

// XMLListAllViolations ...
func getAllViolations(offset, limit int) ([]models.XMLResults, error) {
	var vdetails []models.XMLResults
	findOptions := options.Find()
	findOptions.SetSort(map[string]int{"datetimeunix": -1})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))
	filter := bson.D{}

	cur, err := app.FindDataMongoDB("current_violations", filter, findOptions)
	if err != nil {
		return vdetails, err
	}
	for cur.Next(context.Background()) {
		var dData models.XMLResults
		item := models.CurrentViolations{}
		err := cur.Decode(&item)
		if err != nil {
			continue
		}

		vd := getVehicleDetails(int(item.DeviceID))
		// if vd.SendToNTSA == 0 {
		// 	continue
		// }
		dData.SerialNo = item.DeviceID
		dData.DateOfViolation = item.DateTime.Local().Format("2006-01-02 15:04:05")
		dData.VehicleRegistration = vd.Name
		dData.VehicleOwner = vd.VehicleOwner
		dData.OwnerTel = vd.OwnerTel

		var msg = ""
		if item.Data.Failsafe {
			msg = "Dear Customer, your vehicle " + vd.Name + " had a signal disconnect on " + dData.DateOfViolation
			dData.ViolationType = "Failsafe"
		} else if item.Data.Disconnect {
			msg = "Dear Customer, your vehicle " + vd.Name + " had a power disconnect on " + dData.DateOfViolation
			dData.ViolationType = "Disconnect"
		} else if item.Data.Offline {
			msg = "Dear Customer, your vehicle " + vd.Name + " was offline on " + dData.DateOfViolation
			dData.ViolationType = "Offline"
		} else {
			msg = "Dear Customer, your vehicle " + vd.Name + " was overspeeding on " + dData.DateOfViolation
			dData.ViolationType = "Overspeeding"
		}
		msg += ". Kindly Contact your limiter dealer immediately."

		if dData.VehicleRegistration == "" {
			return vdetails, errors.New("error")
		}

		time.Sleep(3 * time.Second)

		vdetails = append(vdetails, dData)
		sendSMS(vd.VehicleID, item.DeviceID, msg, dData.OwnerTel, dData.ViolationType)
	}

	if err := cur.Err(); err != nil {
		return vdetails, err
	}

	return vdetails, err
}

// GetVehicleDetails ...
func getVehicleDetails(deviceid int) models.VDetails {
	var vd models.VDetails
	query := "SELECT vehicle_details.vehicle_id, send_to_ntsa, vehicle_reg_no, json_value(data, '$.device_detail.owner_name'), json_value(data, '$.device_detail.owner_phone_number') "
	query += " FROM vehicle_configuration "
	query += " LEFT JOIN vehicle_details AS vd ON (vd.vehicle_string_id = vehicle_configuration.vehicle_string_id) "
	query += " WHERE device_id='" + strconv.Itoa(deviceid) + "' LIMIT 1"
	app.DBCon.NewQuery(query).Row(&vd.VehicleID, &vd.SendToNTSA, &vd.Name, &vd.VehicleOwner, &vd.OwnerTel)

	return vd
}
