package updateviolations

import (
	"time"

	"github.com/ekas-portal-api/app"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Status ...
type Status struct {
	// filtered
}

// Run LastSeen.Run() will get triggered automatically.
func (c Status) Run() {
	// select all deviceids
	getAllDeviceIDFromMongoDb()
}

// LastSeen ...
type LastSeen struct {
	ID           int32     `bson:"_id"`
	LastSeenDate time.Time `bson:"last_seen_date"`
	LastSeenUnix uint64    `bson:"last_seen_unix"`
}

func getAllDeviceIDFromMongoDb() {

}

func getVehicleID(deviceID int32) (int32, error) {
	var vid int32
	err := app.DBCon.Select("device_id").From("vehicle_configuration").
		Where(dbx.HashExp{"device_status": "offline"}).Row(&vid)
	return vid, err
}
