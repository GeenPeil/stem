package main

import (
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/GeenPeil/stem/rutte/api"
	"github.com/GeenPeil/stem/rutte/bapi"
	"github.com/GeenPeil/stem/rutte/commonflags"
	"github.com/GeenPeil/stem/rutte/cors"
	"github.com/GeenPeil/stem/rutte/postcodenl"
	"github.com/GeenPeil/stem/rutte/version"

	"github.com/Sirupsen/logrus"
	"github.com/huandu/xstrings"
	flags "github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // package registers driver to database/sql
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	mollieServices "github.com/rollick/gollie/services"
)

// Options type holding all flags/envs for the program
type Options struct {
	commonflags.Log

	Version bool `long:"version" short:"v" description:"Display version and exit"`

	Env string `long:"environment" env:"GPIAC_ENV" default:"local"`

	SelfHTTPAddress string `long:"self-http-address" env:"SELF_HTTP_ADDRESS" default:"https://foobar.geenpeil.nl" description:"HTTP address on which this instance (or a sibbling) can be reached."`

	GPStemDSN          string `long:"gpstem-dsn" env:"GPSTEM_DSN" default:"host='127.0.0.1' dbname=gpstem user=rutte password=rutte sslmode=disable" description:"GeenPeil stem database Postgres Data Source Name"`
	GPStemMaxOpenConns uint   `long:"gpstem-max-open-conns" env:"GPSTEM_MAX_OPEN_CONNS" default:"100" description:"Maximum open connections allowed in the pool for GeenPeil stem database."`

	PostcodeNLAPIUseMock bool   `long:"postcode-nl-api-use-mock" env:"POSTCODE_NL_USE_MOCK" description:"Use mock postcode.nl API implementation"`
	PostcodeNLAPIKey     string `long:"postcode-nl-api-key" env:"POSTCODE_NL_API_KEY" description:"postcode.nl API key"`
	PostcodeNLAPISecret  string `long:"postcode-nl-api-secret" env:"POSTCODE_NL_API_SECRET" description:"postcode.nl API secret"`

	MollieAPIKey string `long:"mollie-api-key" env:"MOLLIE_API_KEY" description:"Mollie PSP API key"`

	HTTPAddress            string `long:"http-addr" env:"HTTP_ADDR" default:":8002" description:"HTTP address to bind to"`
	HTTPEnableWildcardCORS bool   `long:"http-enable-wildcard-cors" description:"Enable HTTP Cross-Origin Resource Sharing"`
}

var options Options

func main() {
	// parse flags
	args, err := flags.Parse(&options)
	if err != nil {
		if et, ok := err.(*flags.Error); ok {
			if et.Type == flags.ErrHelp {
				return
			}
		}
		stdlog.Fatalf("error parsing flags: %v", err)
		return
	}
	if len(args) > 0 {
		stdlog.Fatalf("unexpected arguments: %v", args)
		return
	}

	// validate log level CLI option
	logLevel, err := logrus.ParseLevel(options.Log.Level)
	if err != nil {
		stdlog.Fatalf("invalid log level specified: %v", err)
	}
	if options.Version {
		fmt.Printf("rutte %s\n", version.String())
		os.Exit(0)
	}

	// Setup root logging handler
	logger := logrus.New()
	logger.Level = logLevel
	log := logger.WithFields(logrus.Fields{
		"app": "rutte",
	})
	if options.Log.Format == "json" {
		logger.Formatter = &logrus.JSONFormatter{}
	}
	log.Infoln("starting rutte")

	// Setup DB connection
	if options.GPStemMaxOpenConns == 0 {
		log.Fatal("unlimited maximum open connections is not supported")
	}
	db, err := sqlx.Connect("postgres", options.GPStemDSN)
	if err != nil {
		log.WithError(err).Fatal("error setting up connection to postgres")
	}
	db.MapperFunc(xstrings.ToSnakeCase)

	// Setup Postcode.nl API
	var postcodeAPI postcodenl.API
	if options.PostcodeNLAPIUseMock {
		postcodeAPI = postcodenl.NewMock()
	} else {
		postcodeAPI = postcodenl.New(options.PostcodeNLAPIKey, options.PostcodeNLAPISecret)
	}

	// Setup Mollie payment service API
	molliePaymentService := mollieServices.NewPaymentService(options.MollieAPIKey)

	// Create HTTP server
	logWriter := logger.Writer()
	server := http.Server{
		ErrorLog: stdlog.New(logWriter, "", 0),
	}
	r := chi.NewRouter()
	server.Handler = r
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Route("/api", api.New(log, db, postcodeAPI, molliePaymentService, options.SelfHTTPAddress).AttachChiRouter)
	r.Route("/backoffice-api", bapi.New(log, db).AttachChiRouter)

	if options.Env == "local" {
		aribProxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "localhost:3000"})
		r.NotFound(aribProxy.ServeHTTP)
	}

	// Optionally add CORS headers to each request
	if options.HTTPEnableWildcardCORS {
		log.Info("HTTP wildcard cors headers enabled")
		server.Handler = cors.WrapWildcardCORSHandler(server.Handler)
	} else {
		log.Info("HTTP wildcard cors headers are NOT enabled")
	}

	listener, err := net.Listen("tcp", options.HTTPAddress)
	if err != nil {
		log.WithError(err).WithField("addr", options.HTTPAddress).Fatal("error opening tcp listener")
	}

	log.Infof("starting HTTP server on %s", options.HTTPAddress)
	err = server.Serve(listener)
	if err != nil {
		log.WithError(err).Fatal("error serving HTTP")
	}
}
