package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ekas-portal-api/apis"
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/daos"
	"github.com/ekas-portal-api/errors"
	"github.com/ekas-portal-api/services"
	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
	"github.com/go-ozzo/ozzo-routing/cors"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// load application configurations
	if err := app.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("Invalid application configuration: %s", err))
	}

	// load error messages
	if err := errors.LoadMessages(app.Config.ErrorFile); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	// create the logger
	logger := logrus.New()

	// connect to the database
	dns := getDNS()
	db := app.InitializeDB(dns)
	db.LogFunc = logger.Infof

	err := app.InitializeRedis()
	if err != nil {
		logger.Error(err)
	}

	// wire up API routing
	http.Handle("/", buildRouter(logger, db))

	// start the server
	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	logger.Infof("server %v is started at %v\n", app.Version, address)
	panic(http.ListenAndServe(address, nil))
}

func getDNS() string {
	if os.Getenv("GO_ENV") == "production" {
		return app.Config.ServerDSN
	}

	return app.Config.LocalDSN
}

func buildRouter(logger *logrus.Logger, db *dbx.DB) *routing.Router {
	router := routing.New()

	router.To("GET,HEAD", "/ping", func(c *routing.Context) error {
		c.Abort() // skip all other middlewares/handlers
		return c.Write("OK " + app.Version)
	})

	router.Use(
		app.Init(logger),
		content.TypeNegotiator(content.JSON),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
		app.Transactional(db),
	)

	rg := router.Group("/api/v" + app.Version)

	userDAO := daos.NewUserDAO()
	apis.ServeUserResource(rg, services.NewUserService(userDAO))

	// rg.Post("/auth", apis.Auth(app.Config.JWTSigningKey))
	// rg.Use(auth.JWT(app.Config.JWTVerificationKey, auth.JWTOptions{
	// 	SigningMethod: app.Config.JWTSigningMethod,
	// 	TokenHandler:  apis.JWTHandler,
	// }))

	// artistDAO := daos.NewArtistDAO()
	// apis.ServeArtistResource(rg, services.NewArtistService(artistDAO))

	vehicleDAO := daos.NewVehicleDAO()
	apis.ServeVehicleResource(rg, services.NewVehicleService(vehicleDAO))

	trackingServerServiceDAO := daos.NewTrackingServerServiceDAO()
	apis.ServeTrackingServerServiceResource(rg, services.NewTrackingServerServiceService(trackingServerServiceDAO))

	trackingServerDAO := daos.NewTrackingServerDAO()
	apis.ServeTrackingServerResource(rg, services.NewTrackingServerService(trackingServerDAO))

	deviceDAO := daos.NewDeviceDAO()
	apis.ServeDeviceResource(rg, services.NewDeviceService(deviceDAO))

	vehicleRecordDAO := daos.NewVehicleRecordDAO()
	apis.ServeVehicleRecordResource(rg, services.NewVehicleRecordService(vehicleRecordDAO))

	settingDAO := daos.NewSettingDAO()
	apis.ServeSettingResource(rg, services.NewSettingService(settingDAO))

	return router
}
