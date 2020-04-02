package lastdata

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ekas-portal-api/app"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// LastDataStatus check if device has sent data for the last 30 daya
// then update device ntsa status as false (send_to_ntsa).
type LastDataStatus struct {
	// filtered
}

// Run LastDataStatus.Run() will get triggered automatically.
func (c LastDataStatus) Run() {
	// select all deviceids
	deviceids, err := getAllDeviceIDs()
	if err != nil || len(deviceids) == 0 {
		return
	}

	for _, data := range deviceids {
		//check if table exist
		datetimestamp, err := checkIfDataTableExists(data)
		if datetimestamp == 0 || err != nil {
			continue
		}

		// fmt.Println(datetimestamp)
	}

}

type devices struct {
	DeviceID   int32 `json:"device_id"`
	VehicleID  int32 `json:"vehicle_id"`
	SendToNTSA int32 `json:"send_to_ntsa"`
}

func getAllDeviceIDs() ([]devices, error) {
	// Queries the DB
	// check vehicles without data status (send_to_ntsa)
	deviceids := []devices{}
	err := app.DBCon.Select("vehicle_details.vehicle_id", "device_id", "send_to_ntsa").From("vehicle_details").
		LeftJoin("vehicle_configuration", dbx.NewExp("vehicle_configuration.vehicle_string_id = vehicle_details.vehicle_string_id")).
		Where(dbx.And(dbx.NewExp("device_id>0"), dbx.HashExp{"vehicle_status": 1})).All(&deviceids)
	return deviceids, err
}

// Check if deviceTable Exists
func checkIfDataTableExists(data devices) (int64, error) {
	var exist int
	tx, _ := app.SecondDBCon.Begin()

	query := "SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'ekas_portal_data') AND (TABLE_NAME = 'data_" + strconv.Itoa(int(data.DeviceID)) + "')"
	if err := tx.NewQuery(query).Row(&exist); err != nil {
		return 0, err
	}
	if exist == 0 && data.SendToNTSA > 0 {
		updateNTSAStatus(data.VehicleID, 0)
		return 0, nil
	}

	lastMonth := strconv.FormatInt(getLastMonthUnix(), 10)
	did := strconv.FormatInt(int64(data.DeviceID), 10)
	// Delete data Older tha 30 Days
	query = "DELETE FROM data_" + did + " WHERE date_time_stamp < " + lastMonth
	_, err := app.SecondDBCon.NewQuery(query).Execute()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// select vehicle data current timestamp
	var datetimestamp int64
	if err := app.SecondDBCon.Select("date_time_stamp").From("data_" + did).OrderBy("date_time_stamp DESC").Row(&datetimestamp); err != nil {
		tx.Rollback()
		return 0, err
	}

	fmt.Println(datetimestamp)

	if datetimestamp == 0 && data.SendToNTSA > 0 {
		// update ntsa status as true (send_to_ntsa)
		updateNTSAStatus(data.VehicleID, 0)
	}

	tx.Commit()

	return datetimestamp, nil
}

func getLastMonthUnix() int64 {
	t := time.Now()
	t2 := t.AddDate(0, -1, 0)
	return t2.Unix()
}

func updateNTSAStatus(vid int32, status int) {
	// update ntsa status as false (send_to_ntsa)
	app.DBCon.Update("vehicle_details", dbx.Params{
		"send_to_ntsa": status,
	}, dbx.HashExp{"vehicle_id": vid}).Execute()
}
