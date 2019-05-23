package app

import (
	"fmt"
	"math/rand"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/sony/sonyflake"
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
	LocalDSN  string `mapstructure:"local_dsn"`
	ServerDSN string `mapstructure:"server_dns"`
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
