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
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
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
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	count, err := collection.CountDocuments(ctx, filter, nil)
	return int(count), err
}

// GetDeviceDataLogsMongo ...
func GetDeviceDataLogsMongo(deviceid string, filter primitive.D, opts *options.FindOptions) ([]models.DeviceData, error) {
	var tdetails []models.DeviceData
	// Get collection
	collection := MongoDB.Collection("data_" + deviceid)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

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
	fmt.Println("Found a document: ", strconv.Itoa(i))
	if err := cur.Err(); err != nil {
		return tdetails, err
	}

	return tdetails, err

}
