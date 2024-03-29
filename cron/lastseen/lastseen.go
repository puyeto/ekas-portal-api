package lastseen

import (
	"context"
	"fmt"
	"time"

	"github.com/ekas-portal-api/app"
	dbx "github.com/go-ozzo/ozzo-dbx"
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
	getAllDeviceIDFromMongoDb()
}

// LastSeen ...
type LastSeen struct {
	ID           int32     `bson:"_id"`
	LastSeenDate time.Time `bson:"last_seen_date"`
	LastSeenUnix uint64    `bson:"last_seen_unix"`
}

func getAllDeviceIDFromMongoDb() {
	// Get collection
	collection := app.MongoDB.Collection("a_device_lastseen")
	findOptions := options.Find()
	findOptions.SetSort(map[string]int{"_id": -1})
	findOptions.SetLimit(10000)

	// defer cancel()
	cursor, err := collection.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		return
	}
	// defer cursor.Close(ctx)
	i := 0

	for cursor.Next(context.TODO()) {
		i++
		var m = LastSeen{}
		if err = cursor.Decode(&m); err != nil {
			continue
		}

		deviceStatus := "online"
		// a, _ := time.Parse("2006-01-02 15:04", strDate)
		delta := time.Now().Sub(m.LastSeenDate)

		// update last seen
		if delta.Hours() > 24 && delta.Hours() < (24*14) {
			deviceStatus = "idle"
		} else if delta.Hours() > 72 {
			deviceStatus = "offline"
			// vehicleID, err := getVehicleID(m.ID)
			// if err == nil {
			// 	continue
			// }

			// app.DBCon.Update("vehicle_details", dbx.Params{
			// 	"send_to_ntsa": 0,
			// }, dbx.HashExp{"vehicle_id": vehicleID}).Execute()
		}

		_, err := app.DBCon.Update("vehicle_configuration", dbx.Params{
			"last_seen":     m.LastSeenDate,
			"device_status": deviceStatus,
		}, dbx.HashExp{"device_id": m.ID}).Execute()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("count %v", i)
}

func getVehicleID(deviceID int64) (int32, error) {
	var vid int32
	err := app.DBCon.Select("vehicle_id").From("vehicle_configuration").
		Where(dbx.HashExp{"device_id": deviceID}).Row(&vid)
	return vid, err
}
