package checkdata

import (
	"fmt"
	"strconv"

	"github.com/ekas-portal-api/app"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// CheckDataStatus check if device has sent data every 5 min
// then update device ntsa status as true (send_to_ntsa).
type CheckDataStatus struct {
	// filtered
}

// Run CheckDataStatus.Run() will get triggered automatically.
func (c CheckDataStatus) Run() {
	// select all deviceids
	deviceids, err := getAllDeviceIDs()
	if err != nil || len(deviceids) == 0 {
		return
	}

	for _, data := range deviceids {
		// check if table exist
		exist, err := checkIfDataTableExists(data.DeviceID)
		if exist == 0 || err != nil {
			continue
		}

		// Check if device has sent data last 3 days.
		confirm, err := confirmPreviousData(3)
		if confirm == 0 || err != nil {
			continue
		}

		// update ntsa status as true (send_to_ntsa)
		// app.DBCon.Update("vehicle_details", dbx.Params{
		// 	"send_to_ntsa": 1,
		// }, dbx.HashExp{"key_string": data.VehicleID}).Execute()
		fmt.Println(data.DeviceID, data.VehicleID, exist)

	}
}

type devices struct {
	DeviceID  int32 `json:"device_id"`
	VehicleID int32 `json:"vehicle_id"`
}

func getAllDeviceIDs() ([]devices, error) {
	// Queries the DB
	// check vehicles without data status (send_to_ntsa)
	deviceids := []devices{}
	err := app.DBCon.Select("vehicle_details.vehicle_id", "device_id").From("vehicle_details").
		LeftJoin("vehicle_configuration", dbx.NewExp("vehicle_configuration.vehicle_string_id = vehicle_details.vehicle_string_id")).
		Where(dbx.HashExp{"send_to_ntsa": 0}).All(&deviceids)
	return deviceids, err
}

// Check if deviceTable Exists
func checkIfDataTableExists(id int32) (int, error) {
	var exist int
	query := "SELECT count(*) FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'ekas_portal_data') AND (TABLE_NAME = 'data_" + strconv.Itoa(int(id)) + "')"
	err := app.SecondDBCon.NewQuery(query).Row(&exist)

	return exist, err
}

func confirmPreviousData(noofdays int) (int, error) {
	return 0, nil
}
