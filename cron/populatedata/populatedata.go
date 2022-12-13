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
	data, _ := getTripDataByDeviceIDBtwDates("1921225606", 1669923933, 1770040060, 151, 100)
	// var previous time.Time
	fmt.Println(len(data))
	for i := len(data) - 1; i >= 0; i-- {
		// if previous != data[i].DateTime {
		data[i].DeviceID = 1921225605
		// data[i].GroundSpeed = 0.00
		// if data[i].Latitude < -5000000 {
		// 	data[i].Latitude = data[i].Latitude + 4000000
		// } else if data[i].Latitude > -4000000 {
		// 	data[i].Latitude = data[i].Latitude + 2000000
		// }

		data[i].UTCTimeDay = data[i].UTCTimeDay - 2
		data[i].UTCTimeHours = data[i].UTCTimeHours - 1

		dt := data[i].DateTime

		data[i].DateTime = dt.AddDate(0, 0, -2)
		data[i].DeviceTime = dt.AddDate(0, 0, -2)

		data[i].DateTime = data[i].DateTime.Add(-time.Hour * 1)
		data[i].DeviceTime = data[i].DeviceTime.Add(-time.Hour * 1)

		if data[i].UTCTimeMinutes > 30 {
			data[i].UTCTimeMinutes = data[i].UTCTimeMinutes - 30
			data[i].DateTime = data[i].DateTime.Add(-time.Minute * 30)
			data[i].DeviceTime = data[i].DeviceTime.Add(-time.Minute * 30)
		}
		data[i].DateTimeStamp = data[i].DateTime.Unix()

		fmt.Println(data[i].DateTime, data[i].DateTimeStamp, data[i].UTCTimeHours, data[i].UTCTimeMinutes)
		LogToRedis(data[i])
		app.LogToMongoDB(data[i])
		app.LoglastSeenMongoDB(data[i])
		// previous = data[i].DateTime
		// }
		if i == 0 {
			fmt.Println("Finished")
		}
	}
}

func getTripDataByDeviceIDBtwDates(deviceid string, from, to int64, offset, limit int) ([]models.DeviceData, error) {
	fmt.Println("Job Started")
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
	data, _ := deleteTripDataByDeviceIDBtwDates(2023202064, 1610256527, 1610472527)
	log.Println(data)
}

func deleteTripDataByDeviceIDBtwDates(id, from, to uint32) (int, error) {
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
	if err := app.SetValue(key, data); err != nil {
		fmt.Println(err)
	}
}
