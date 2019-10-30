package app

import (
	"encoding/json"
	"time"

	"github.com/ekas-portal-api/models"
	"github.com/go-redis/redis"
)

var (
	redisClient *redis.Client
	// redisURL    = "db-redis-cluster-do-user-4666162-0.db.ondigitalocean.com:25061"

)

// InitializeRedis ...
func InitializeRedis() error {
	// if os.Getenv("GO_ENV") != "production" {
	// 	redisURL = "db-redis-cluster-do-user-4666162-0.db.ondigitalocean.com:25061"
	// }

	opt, _ := redis.ParseURL("rediss://default:wdbsxehbizfl5kbu@db-redis-cluster-do-user-4666162-0.db.ondigitalocean.com:25061/1")
	opt.PoolSize = 100
	opt.MaxRetries = 2
	opt.ReadTimeout = -1

	redisClient = redis.NewClient(opt)

	ping, err := redisClient.Ping().Result()
	if err == nil && len(ping) > 0 {
		println("Connected to Redis")
		return nil
	}
	println("Redis Connection Failed")
	return err
}

// GetValue ...
func GetValue(key string) (interface{}, error) {
	var deserializedValue interface{}
	serializedValue, err := redisClient.Get(key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue, err
}

// GetDeviceDataValue ...
func GetDeviceDataValue(key string) (models.DeviceData, error) {
	var deserializedValue models.DeviceData
	serializedValue, err := redisClient.Get(key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue, err
}

// GetLastSeenValue ...
func GetLastSeenValue(key string) (models.DeviceData, error) {
	var deserializedValue models.LastSeenStruct
	serializedValue, err := redisClient.Get(key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue.DeviceData, err
}

// SetValue ...
func SetValue(key string, value interface{}) (bool, error) {
	serializedValue, _ := json.Marshal(value)
	err := redisClient.Set(key, string(serializedValue), 0).Err()
	return true, err
}

// SetValueWithTTL ...
func SetValueWithTTL(key string, value interface{}, ttl int) (bool, error) {
	serializedValue, _ := json.Marshal(value)
	err := redisClient.Set(key, string(serializedValue), time.Duration(ttl)*time.Second).Err()
	return true, err
}

// RPush ...
func RPush(key string, valueList []string) (bool, error) {
	err := redisClient.RPush(key, valueList).Err()
	return true, err
}

// RpushWithTTL ...
func RpushWithTTL(key string, valueList []string, ttl int) (bool, error) {
	err := redisClient.RPush(key, valueList, ttl).Err()
	return true, err
}

// LRange ...
func LRange(key string, start, stop int64) ([]string, error) {
	val, err := redisClient.LRange(key, start, stop).Result()
	return val, err
}

// ZRange ...
// Returns the specified range of elements in the sorted set stored at key
func ZRange(key string, start, stop int64) ([]string, error) {
	val, err := redisClient.ZRange(key, start, stop).Result()
	return val, err
}

// ZRevRange ...
// Returns the specified range of elements in the sorted set stored at key
// elements are considered to be ordered from high to low scores
func ZRevRange(key string, start, stop int64) ([]string, error) {
	val, err := redisClient.ZRevRange(key, start, stop).Result()
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
	val, err := redisClient.ZRevRangeByScore(key, &opt).Result()
	return val, err
}

// ZCount ...
func ZCount(key, min, max string) int64 {
	return redisClient.ZCount(key, min, max).Val()
}

// ListLength ...
func ListLength(key string) int64 {
	return redisClient.LLen(key).Val()
}

// Publish ...
func Publish(channel string, message string) {
	redisClient.Publish(channel, message)
}

// GetKeyListByPattern ...
func GetKeyListByPattern(pattern string) []string {
	return redisClient.Keys(pattern).Val()
}

// IncrementValue ...
func IncrementValue(key string) int64 {
	return redisClient.Incr(key).Val()
}

// DelKey ...
func DelKey(key string) error {
	return redisClient.Del(key).Err()
}

// ListKeys ...
func ListKeys(key string) ([]string, error) {
	return redisClient.Keys(key).Result()
}
