package app

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rs/xid"
	"github.com/sony/sonyflake"
	"github.com/speps/go-hashids"
	"github.com/spf13/viper"
)

// Config stores the application-wide configurations
var Config appConfig

type appConfig struct {
	// the path to the error message file. Defaults to "config/errors.yaml"
	ErrorFile string `mapstructure:"error_file"`
	// the server port. Defaults to 8080
	ServerPort int `mapstructure:"server_port"`
	// the data source name (DSN) for connecting to the database. required.
	LocalDSN    string `mapstructure:"local_dsn"`
	ServerDSN   string `mapstructure:"server_dns"`
	MongoDBDNS  string `mapstructure:"mongo_db_dns"`
	MongoDBName string `mapstructure:"mongo_db_name"`

	// the signing method for JWT. Defaults to "HS256"
	JWTSigningMethod string `mapstructure:"jwt_signing_method"`
	// JWT signing key. required.
	JWTSigningKey string `mapstructure:"jwt_signing_key"`
	// JWT verification key. required.
	JWTVerificationKey string `mapstructure:"jwt_verification_key"`
	// url to tracking server
	TrackingServerURL string `mapstructure:"tracking_server_url"`
	// Twilio
	TwilioAccountSID string `mapstructure:"twilio_account_sid"`
	TwilioAuthToken  string `mapstructure:"twilio_auth_token"`
}

func (config appConfig) Validate() error {
	return validation.ValidateStruct(&config,
		validation.Field(&config.LocalDSN, validation.Required),
		validation.Field(&config.ServerDSN, validation.Required),
		validation.Field(&config.JWTSigningKey, validation.Required),
		validation.Field(&config.JWTVerificationKey, validation.Required),
	)
}

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
// Environment variables with the prefix "RESTFUL_" in their names are also read automatically.
func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("restful")
	v.AutomaticEnv()
	v.SetDefault("error_file", "config/errors.yaml")
	v.SetDefault("server_port", 8081)
	v.SetDefault("jwt_signing_method", "HS256")
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("Failed to read the configuration file: %s", err)
	}

	if err := v.Unmarshal(&Config); err != nil {
		return err
	}
	return Config.Validate()
}

// GenerateNewID1 Generate new id using
// Note: this is base16, could shorten by encoding as base62 string
// fmt.Printf("github.com/sony/sonyflake:   %x\n", id)
func GenerateNewID1() uint64 {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, _ := flake.NextID()
	return id
}

// GenerateNewID ...
func GenerateNewID() uint32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Uint32()
}

// GenerateNewStringID ...
func GenerateNewStringID() string {
	guid := xid.New()
	hd := hashids.NewData()
	hd.Salt = guid.String()
	hd.MinLength = 8
	hd.Alphabet = "ABCDEFGHIJKLMNPQRSTUVWXYZ123456789"
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode([]int{45, 434, 1313})

	return e
}

// CalculatePassHash calculate password hash usin paswword and salt
func CalculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// RandStringBytes ...
func RandStringBytes(letterBytes string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
