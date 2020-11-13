package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bamzi/jobrunner"
	"github.com/ekas-portal-api/apis"
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/cron/lastseen"
	"github.com/ekas-portal-api/cron/populatedata"
	"github.com/ekas-portal-api/cron/updateviolations"
	"github.com/ekas-portal-api/daos"
	"github.com/ekas-portal-api/errors"
	"github.com/ekas-portal-api/services"
	dbx "github.com/go-ozzo/ozzo-dbx"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
	"github.com/go-ozzo/ozzo-routing/v2/content"
	"github.com/go-ozzo/ozzo-routing/v2/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
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
	app.InitLogger(logger)

	// connect to the database
	dns := getDNS()
	db := app.InitializeDB(dns)
	db.LogFunc = logger.Infof

	// Connect to mongodb
	app.MongoDB = app.InitializeMongoDB(app.Config.MongoDBDNS, app.Config.MongoDBName, logger)

	err := app.InitializeRedis()
	if err != nil {
		logger.Error(err)
	}

	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	// jobrunner.Schedule("@every 1m", checkexpired.Status{})

	if os.Getenv("GO_ENV") == "production" {
		// run cronjobs
		// jobrunner.Schedule("CRON_TZ=Africa/Nairobi * 8 * * *", checkexpired.Status{})
		// go jobrunner.Schedule("@every 60m", checkdata.Status{})
		jobrunner.Schedule("@every 60m", lastseen.Status{})
		// go jobrunner.In(10*time.Second, updateviolations.Status{})
		jobrunner.Schedule("@midnight", updateviolations.Status{}) // every midnight do this..
		jobrunner.Schedule("@every 60m", lastseen.Status{})
	} else {
		jobrunner.In(2*time.Second, populatedata.Status{})
	}

	// wire up API routing
	http.Handle("/", buildRouter(logger, db))

	// start the server
	address := fmt.Sprintf(":%v", app.Config.ServerPort)
	httsaddress := fmt.Sprintf(":%v", app.Config.ServerPort+1)
	logger.Infof("server %v is started at %v (http) and %v (https)\n", app.Version, address, httsaddress)
	// panic(http.ListenAndServeTLS(address, "server.rsa.crt", "server.rsa.key", nil))

	//  Start HTTP
	go func() {
		panic(http.ListenAndServe(address, nil))
	}()

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		// ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8082 with the TLS config
	server := &http.Server{
		Addr: ":8082",
		// TLSConfig: tlsConfig,
		Handler: buildRouter(logger, db),
	}

	// Listen to HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))

	//  Start HTTPS
	// panic(server.ListenAndServeTLS("./server.rsa.crt", "./server.rsa.key"))
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

	// artistDAO := daos.NewArtistDAO()
	// apis.ServeArtistResource(rg, services.NewArtistService(artistDAO))

	vehicleDAO := daos.NewVehicleDAO()
	apis.ServeVehicleResource(rg, services.NewVehicleService(vehicleDAO))

	trackingServerServiceDAO := daos.NewTrackingServerServiceDAO()
	apis.ServeTrackingServerServiceResource(rg, services.NewTrackingServerServiceService(trackingServerServiceDAO))

	trackingServerDAO := daos.NewTrackingServerDAO()
	apis.ServeTrackingServerResource(rg, services.NewTrackingServerService(trackingServerDAO))

	settingDAO := daos.NewSettingDAO()
	apis.ServeSettingResource(rg, services.NewSettingService(settingDAO))

	companyDAO := daos.NewCompanyDAO()
	apis.ServeCompanyResource(rg, services.NewCompanyService(companyDAO))

	if os.Getenv("GO_ENV") == "production" {
		rg.Use(auth.JWT(app.Config.JWTVerificationKey, auth.JWTOptions{
			SigningMethod: app.Config.JWTSigningMethod,
			TokenHandler:  app.JWTHandler,
		}))
	}

	deviceDAO := daos.NewDeviceDAO()
	apis.ServeDeviceResource(rg, services.NewDeviceService(deviceDAO))

	vehicleRecordDAO := daos.NewVehicleRecordDAO()
	apis.ServeVehicleRecordResource(rg, services.NewVehicleRecordService(vehicleRecordDAO))

	certificateDAO := daos.NewCertificateDAO()
	apis.ServeCertificateResource(rg, services.NewCertificateService(certificateDAO))

	discountDAO := daos.NewDiscountDAO()
	apis.ServeDiscountResource(rg, services.NewDiscountService(discountDAO))

	pricingDAO := daos.NewPricingDAO()
	apis.ServePricingResource(rg, services.NewPricingService(pricingDAO))

	invoiceDAO := daos.NewInvoiceDAO()
	apis.ServeInvoiceResource(rg, services.NewInvoiceService(invoiceDAO))

	ownersDAO := daos.NewOwnerDAO()
	apis.ServeOwnerResource(rg, services.NewOwnerService(ownersDAO))

	simcardsDAO := daos.NewSimcardDAO()
	apis.ServeSimcardResource(rg, services.NewSimcardService(simcardsDAO))

	return router
}
