package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ekas-portal-api/models"
	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
	// redisURL    = "db-redis-cluster-do-user-4666162-0.db.ondigitalocean.com:25061"

)

// InitializeRedis ...
func InitializeRedis() error {
	// new redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "api.equscabanus.com:6379",
		Password: "",
		DB:       3,
	})

	// test connection
	ping, err := redisClient.Ping(ctx).Result()
	if err == nil && len(ping) > 0 {
		// FlushAll()
		println("Connected to Redis")
		return nil
	}
	println("Redis Connection Failed")
	return err
}

// GetValue ...
func GetValue(key string) (interface{}, error) {
	var deserializedValue interface{}
	serializedValue, err := redisClient.Get(ctx, key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue, err
}

// GetDeviceDataValue ...
func GetDeviceDataValue(key string) (models.DeviceData, error) {
	var deserializedValue models.DeviceData
	serializedValue, err := redisClient.Get(ctx, key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue, err
}

// GetLastSeenValue ...
func GetLastSeenValue(key string) (models.DeviceData, error) {
	var deserializedValue models.LastSeenStruct
	serializedValue, err := redisClient.Get(ctx, key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue.DeviceData, err
}

// SetValue ...
func SetValue(key string, value interface{}) error {
	serializedValue, _ := json.Marshal(value)
	return redisClient.Set(ctx, key, string(serializedValue), 0).Err()
}

// SetValueWithTTL ...
func SetValueWithTTL(key string, value interface{}, ttl int) (bool, error) {
	serializedValue, _ := json.Marshal(value)
	err := redisClient.Set(ctx, key, string(serializedValue), time.Duration(ttl)*time.Second).Err()
	return true, err
}

// RPush ...
func RPush(key string, valueList []string) (bool, error) {
	err := redisClient.RPush(ctx, key, valueList).Err()
	return true, err
}

// RpushWithTTL ...
func RpushWithTTL(key string, valueList []string, ttl int) (bool, error) {
	err := redisClient.RPush(ctx, key, valueList, ttl).Err()
	return true, err
}

// LRange ...
func LRange(key string, start, stop int64) ([]string, error) {
	val, err := redisClient.LRange(ctx, key, start, stop).Result()
	return val, err
}

// ZAdd ...
// Adds all the specified members with the specified scores to the sorted set stored at key
func ZAdd(key string, score int64, members interface{}) error {
	serializedValue, _ := json.Marshal(members)
	err := redisClient.ZAdd(ctx, key, &redis.Z{
		Score:  float64(score),
		Member: string(serializedValue),
	}).Err()
	return err
}

// ZRange ...
// Returns the specified range of elements in the sorted set stored at key
// The elements are considered to be ordered from the lowest to the highest
func ZRange(key string, start, stop int64) ([]string, error) {
	val, err := redisClient.ZRange(ctx, key, start, stop).Result()
	return val, err
}

// ZRevRange ...
// Returns the specified range of elements in the sorted set stored at key
// elements are considered to be ordered from high to low scores
func ZRevRange(key string, start, stop int64) ([]string, error) {
	val, err := redisClient.ZRevRange(ctx, key, start, stop).Result()
	return val, err
}

// ZRevRangeByScore ...
// Returns the specified range of elements in the sorted set stored at key
// elements are considered to be ordered from high to low scores
func ZRevRangeByScore(key string, min, max string, offset, lim int64) ([]string, error) {
	opt := redis.ZRangeBy{
		Min:    min,
		Max:    max,
		Offset: offset,
		Count:  lim,
	}
	val, err := redisClient.ZRevRangeByScore(ctx, key, &opt).Result()
	return val, err
}

// ZCount ...
func ZCount(key, min, max string) int64 {
	return redisClient.ZCount(ctx, key, min, max).Val()
}

// ListLength ...
func ListLength(key string) int64 {
	return redisClient.LLen(ctx, key).Val()
}

// Publish ...
func Publish(channel string, message string) {
	redisClient.Publish(ctx, channel, message)
}

// GetKeyListByPattern ...
func GetKeyListByPattern(pattern string) []string {
	return redisClient.Keys(ctx, pattern).Val()
}

// IncrementValue ...
func IncrementValue(key string) int64 {
	return redisClient.Incr(ctx, key).Val()
}

// DelKey ...
func DelKey(key string) error {
	return redisClient.Del(ctx, key).Err()
}

// FlushAll ...
func FlushAll() error {
	return redisClient.FlushAll(ctx).Err()
}

// ListKeys ...
func ListKeys(key string) ([]string, error) {
	return redisClient.Keys(ctx, key).Result()
}

// SetRedisLog log to redis
func SetRedisLog(devData chan models.DeviceData, key string) error {
	dev := <-devData
	err := ZAdd(key, dev.DateTimeStamp, dev)
	fmt.Println(err, key, dev.DateTimeStamp)
	fmt.Println(dev)
	return err
}
