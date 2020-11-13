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
	populatedata()
	// deletedata()
}

func populatedata() {
	// Select data between dates 1166676296
	data, _ := getTripDataByDeviceIDBtwDates("1114211591", 1604640200, 1604850200, 8686, 10)
	var previous time.Time
	log.Println(len(data))
	for i := 0; i < len(data); i++ {
		if previous != data[i].DateTime {
			data[i].DeviceID = 1829209633
			fmt.Println(i, len(data), data[i].DateTime, data[i].Latitude, data[i].Longitude, data[i].GroundSpeed)
			// app.LogToMongoDB(data[i])
			// app.LoglastSeenMongoDB(data[i])
			// previous = data[i].DateTime
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
	data, _ := filterTripDataByDeviceIDBtwDates(1829209633, 1604886200, 1604887200)
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
