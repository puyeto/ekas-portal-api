package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/bamzi/jobrunner"
	"github.com/ekas-portal-api/apis"
	"github.com/ekas-portal-api/app"
	"github.com/ekas-portal-api/cron/checkdata"
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

	// connect to second database
	seconddns := getDNSForSecondDB()
	seconddb := app.InitializeSecondDB(seconddns)
	seconddb.LogFunc = logger.Infof

	err := app.InitializeRedis()
	if err != nil {
		logger.Error(err)
	}

	// run cronjobs
	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	go jobrunner.Schedule("@every 20m", checkdata.CheckDataStatus{})

	// wire up API routing
	http.Handle("/", buildRouter(logger, db, seconddb))

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

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr: ":8082",
		// TLSConfig: tlsConfig,
		Handler: buildRouter(logger, db, seconddb),
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

func getDNSForSecondDB() string {
	if os.Getenv("GO_ENV") == "production" {
		return app.Config.SecondServerDSN
	}

	return app.Config.SecondLocalDSN
}

func buildRouter(logger *logrus.Logger, db *dbx.DB, seconddb *dbx.DB) *routing.Router {
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

	certificateDAO := daos.NewCertificateDAO()
	apis.ServeCertificateResource(rg, services.NewCertificateService(certificateDAO))

	discountDAO := daos.NewDiscountDAO()
	apis.ServeDiscountResource(rg, services.NewDiscountService(discountDAO))

	pricingDAO := daos.NewPricingDAO()
	apis.ServePricingResource(rg, services.NewPricingService(pricingDAO))

	invoiceDAO := daos.NewInvoiceDAO()
	apis.ServeInvoiceResource(rg, services.NewInvoiceService(invoiceDAO))

	companyDAO := daos.NewCompanyDAO()
	apis.ServeCompanyResource(rg, services.NewCompanyService(companyDAO))

	ownersDAO := daos.NewOwnerDAO()
	apis.ServeOwnerResource(rg, services.NewOwnerService(ownersDAO))

	return router
}
