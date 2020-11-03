package updateviolations

import (
	"context"
	"strconv"
	"time"

	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/models"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Status ...
type Status struct {
	// filtered
}

// Run LastSeen.Run() will get triggered automatically.
func (c Status) Run() {
	// select all deviceids
	getAllOfflines()
}

// Devices ...
type Devices struct {
	DeviceID int32     `json:"device_id"`
	LastSeen time.Time `json:"last_seen"`
}

func getAllOfflines() {
	devices, _ := getVehicleID()
	for _, dev := range devices {
		// filter := bson.D{{Key: "deviceid", Value: int(dev.DeviceID)}}
		count, err := Count(strconv.Itoa(int(dev.DeviceID)), bson.D{}, nil)
		if err != nil {
			continue
		}

		if count > 0 && count <= 5 {
			// Delete collection
			collection := app.MongoDB.Collection("data_" + strconv.Itoa(int(dev.DeviceID)))
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			collection.DeleteMany(ctx, bson.D{})
		}
		var m models.DeviceData
		m.DeviceID = uint32(dev.DeviceID)
		m.Offline = true
		diff := time.Now().Sub(dev.LastSeen)
		if diff.Hours() > 168 {
			m.Disconnect = true
		}
		m.TransmissionReason = 255
		m.DateTime = dev.LastSeen
		m.DateTimeStamp = dev.LastSeen.Unix()
		if err := LogCurrentViolationSeenMongoDB(m); err != nil {
			continue
		}
		if err := LogToMongoDB(m); err != nil {
			continue
		}
	}
}

// LogToMongoDB ...
func LogToMongoDB(m models.DeviceData) error {
	collection := app.MongoDB.Collection("data_" + strconv.FormatInt(int64(m.DeviceID), 10))
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_, err := collection.InsertOne(ctx, m)
	return err
}

// Count returns the number of trip records in the database.
func Count(deviceid string, filter primitive.D, opts *options.FindOptions) (int, error) {
	app.CreateIndexMongo("data_" + deviceid)
	collection := app.MongoDB.Collection("data_" + deviceid)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	count, err := collection.CountDocuments(ctx, filter, nil)
	return int(count), err
}

func getVehicleID() ([]Devices, error) {
	devices := []Devices{}
	err := app.DBCon.Select("device_id", "COALESCE(vd.created_on, last_seen) AS last_seen").From("vehicle_configuration AS vc").
		LeftJoin("vehicle_details AS vd", dbx.NewExp("vd.vehicle_id = vc.vehicle_id")).
		Where(dbx.HashExp{"device_status": "offline"}).All(&devices)
	return devices, err
}

// LogCurrentViolationSeenMongoDB update current violation
func LogCurrentViolationSeenMongoDB(m models.DeviceData) error {
	data := bson.M{
		"$set": bson.M{
			"data":         m,
			"datetime":     m.DateTime,
			"datetimeunix": m.DateTimeStamp,
		},
	}

	return upsert(data, m.DeviceID, "current_violations")
}

func upsert(data bson.M, deviceID uint32, table string) error {
	collection := app.MongoDB.Collection(table)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": deviceID}, data, opts)

	return err
}
