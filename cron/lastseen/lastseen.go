package lastseen

import (
	"context"
	"time"

	"github.com/ekas-portal-api/app"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"go.mongodb.org/mongo-driver/bson"
)

type Status struct {
	// filtered
}

// Run LastSeen.Run() will get triggered automatically.
func (c Status) Run() {
	// select all deviceids
	getAllDeviceIDFromMongoDb()
}

type LastSeen struct {
	ID           int32     `bson:"_id"`
	LastSeenDate time.Time `bson:"last_seen_date"`
	LastSeenUnix uint64    `bson:"last_seen_unix"`
}

func getAllDeviceIDFromMongoDb() {
	// Get collection
	collection := app.MongoDB.Collection("a_device_lastseen")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var m = LastSeen{}
		if err = cursor.Decode(&m); err != nil {
			return
		}

		// update last seen
		app.DBCon.Update("vehicle_configuration", dbx.Params{
			"last_seen": m.LastSeenUnix,
		}, dbx.HashExp{"device_id": m.ID}).Execute()
	}
}
