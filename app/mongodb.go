package app

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/ekas-portal-api/models"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB ...
var MongoDB *mongo.Database

// InitializeMongoDB Initialize MongoDB Connection
func InitializeMongoDB(dbURL, dbName string, logger *logrus.Logger) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// defer client.Disconnect(ctx)

	logger.Infof("Mongo DB initialized: %v", dbName)
	return client.Database(dbName)
}

// CountRecordsMongo returns the number of records in the database.
func CountRecordsMongo(colName string, filter primitive.M, opts *options.FindOptions) (int, error) {
	collection := MongoDB.Collection(colName)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter, nil)
	return int(count), err
}

// GetDeviceDataLogsMongo ...
func GetDeviceDataLogsMongo(deviceid string, filter primitive.D, opts *options.FindOptions) ([]models.DeviceData, error) {
	CreateIndexMongo(deviceid)
	var tdetails []models.DeviceData
	// Get collection
	collection := MongoDB.Collection("data_" + deviceid)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return tdetails, err
	}
	defer cur.Close(ctx)
	i := 0

	for cur.Next(context.Background()) {

		item := models.DeviceData{}
		err := cur.Decode(&item)
		if err != nil {
			return tdetails, err
		}
		tdetails = append(tdetails, item)
		i++
	}
	// fmt.Println("Found a document: ", strconv.Itoa(i))
	if err := cur.Err(); err != nil {
		return tdetails, err
	}

	return tdetails, err

}

// FindDataMongoDB ...
func FindDataMongoDB(colname string, filter primitive.D, opts *options.FindOptions) (*mongo.Cursor, error) {

	// Get collection
	collection := MongoDB.Collection(colname)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.Find(ctx, filter, opts)
}

// CreateIndexMongo create a mongodn index
func CreateIndexMongo(deviceid string) (string, error) {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"datetimestamp": -1, // index in ascending order
		}, Options: nil,
	}
	collection := MongoDB.Collection("data_" + deviceid)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return collection.Indexes().CreateOne(ctx, mod)
}

// LogToMongoDB ...
func LogToMongoDB(m models.DeviceData) error {
	fmt.Println(m)
	collection := MongoDB.Collection("data_" + strconv.FormatInt(int64(m.DeviceID), 10))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, m)
	CreateIndexMongo("data_" + strconv.FormatInt(int64(m.DeviceID), 10))
	return err
}

// LoglastSeenMongoDB update last seen
func LoglastSeenMongoDB(m models.DeviceData) error {
	data := bson.M{
		"$set": bson.M{
			"last_seen_date": m.DateTime,
			"last_seen_unix": m.DateTimeStamp,
			"updated_at":     time.Now(),
		},
	}

	return upsert(data, m.DeviceID, "a_device_lastseen")
}

func upsert(data bson.M, deviceID uint64, table string) error {
	collection := MongoDB.Collection(table)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, bson.M{"_id": deviceID}, data, opts)

	return err
}
