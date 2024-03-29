package checkdata

import (
	"strconv"

	"github.com/ekas-portal-api/app"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"go.mongodb.org/mongo-driver/bson"
)

// Status check if device has sent data every 5 min
// then update device ntsa status as true (send_to_ntsa).
type Status struct {
	// filtered
}

// Run Status.Run() will get triggered automatically.
func (c Status) Run() {
	// select all deviceids
	deviceids, err := getAllDeviceIDs()
	if err != nil || len(deviceids) == 0 {
		return
	}

	for _, data := range deviceids {
		// check if table exist
		filter := bson.M{"deviceid": data.DeviceID}
		count, err := app.CountRecordsMongo("data_"+strconv.Itoa(int(data.DeviceID)), filter, nil)
		if count == 0 || err != nil {
			continue
		}

		// Check if device has sent data last 3 days.
		// confirm, err := confirmPreviousData(3)
		// if confirm == 0 || err != nil {
		// 	continue
		// }

		// update ntsa status as true (send_to_ntsa)
		app.DBCon.Update("vehicle_details", dbx.Params{
			"send_to_ntsa": 1,
		}, dbx.HashExp{"vehicle_id": data.VehicleID}).Execute()

	}
}

type devices struct {
	DeviceID  int64 `json:"device_id"`
	VehicleID int32 `json:"vehicle_id"`
}

func getAllDeviceIDs() ([]devices, error) {
	// Queries the DB
	// check vehicles without data status (send_to_ntsa)
	deviceids := []devices{}
	err := app.DBCon.Select("vehicle_details.vehicle_id", "device_id").From("vehicle_details").
		LeftJoin("vehicle_configuration", dbx.NewExp("vehicle_configuration.vehicle_string_id = vehicle_details.vehicle_string_id")).
		Where(dbx.HashExp{"send_to_ntsa": 0, "vehicle_status": 1, "speed_source": 1}).Limit(30).OrderBy("RAND()").All(&deviceids)
	return deviceids, err
}

func confirmPreviousData(noofdays int) (int, error) {
	return 0, nil
}
