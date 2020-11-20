package populatedata

import (
	"context"
	"fmt"
	"log"
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
	// go deletedata()
	populatedata()
}

func populatedata() {
	// Select data between dates
	data, _ := getTripDataByDeviceIDBtwDates("1728208958", 1599220800, 1604491200, 10000, 5000)
	var previous time.Time
	for i := 0; i < len(data); i++ {
		if previous != data[i].DateTime {
			data[i].DeviceID = 2004205614
			// data[i].GroundSpeed = 0.00
			// if data[i].Latitude < -5000000 {
			// 	data[i].Latitude = data[i].Latitude + 4000000
			// } else if data[i].Latitude > -4000000 {
			// 	data[i].Latitude = data[i].Latitude + 2000000
			// }
			fmt.Println(i, len(data), data[i].DateTime, data[i].Latitude, data[i].Longitude, data[i].GroundSpeed)
			LogToRedis(data[i])
			app.LogToMongoDB(data[i])
			app.LoglastSeenMongoDB(data[i])
			previous = data[i].DateTime
		}
	}
}

func getTripDataByDeviceIDBtwDates(deviceid string, from, to int64, offset, limit int) ([]models.DeviceData, error) {
	findOptions := options.Find()
	// Sort by `price` field descending
	findOptions.SetSort(map[string]int{"datetimestamp": 1})
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	filter := bson.D{
		{Key: "datetimestamp", Value: bson.D{{Key: "$gte", Value: from}}},
		{Key: "datetimestamp", Value: bson.D{{Key: "$lte", Value: to}}},
	}
	return app.GetDeviceDataLogsMongo(deviceid, filter, findOptions)
}

func deletedata() {
	// Select data between dates
	data, _ := filterTripDataByDeviceIDBtwDates(1829209633, 1604722200, 1605887200)
	log.Println(data)
}

func filterTripDataByDeviceIDBtwDates(id, from, to uint32) (int, error) {
	filter := bson.D{
		{Key: "datetimestamp", Value: bson.D{{Key: "$gte", Value: from}}},
		{Key: "datetimestamp", Value: bson.D{{Key: "$lte", Value: to}}},
	}

	// Get collection
	collection := app.MongoDB.Collection("data_" + strconv.Itoa(int(id)))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(res.DeletedCount), nil
}

// LogToRedis log data to redis
func LogToRedis(m models.DeviceData) {
	var device = strconv.FormatUint(uint64(m.DeviceID), 10)
	lastSeen(m, "lastseen:"+device)
	lastSeen(m, "lastseen")
	// if m.TransmissionReason != 255 && m.GroundSpeed != 0 {
	SetRedisLog(m, "data:"+device)
	// }
}

// SetRedisLog log to redis
func SetRedisLog(m models.DeviceData, key string) {
	err := app.ZAdd(key, m.DateTimeStamp, m)
	if err != nil {
		fmt.Println(err)
	}
}

type lastSeenStruct struct {
	DateTime   time.Time
	DeviceData models.DeviceData
}

func lastSeen(m models.DeviceData, key string) {
	var data = lastSeenStruct{
		DateTime:   m.DateTime,
		DeviceData: m,
	}
	// SET object
	_, err := app.SetValue(key, data)
	if err != nil {
		fmt.Println(err)
	}
}
