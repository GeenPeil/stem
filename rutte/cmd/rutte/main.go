package main

import (
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"os"

	"github.com/GeenPeil/stem/rutte/api"
	"github.com/GeenPeil/stem/rutte/bapi"
	"github.com/GeenPeil/stem/rutte/commonflags"
	"github.com/GeenPeil/stem/rutte/cors"
	"github.com/GeenPeil/stem/rutte/version"

	"github.com/Sirupsen/logrus"
	"github.com/huandu/xstrings"
	flags "github.com/jessevdk/go-flags"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // package registers driver to database/sql
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

// Options type holding all flags/envs for the program
type Options struct {
	commonflags.Log

	Version bool `long:"version" short:"v" description:"Display version and exit"`

	GPStemDSN          string `long:"gpstem-dsn" env:"GPSTEM_DSN" default:"host='127.0.0.1' dbname=gpstem user=rutte password=rutte sslmode=disable" description:"GeenPeil stem database Postgres Data Source Name"`
	GPStemMaxOpenConns uint   `long:"gpstem-max-open-conns" env:"GPSTEM_MAX_OPEN_CONNS" default:"100" description:"Maximum open connections allowed in the pool for GeenPeil stem database."`

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
	r.Route("/api", api.New(log, db).AttachChiRouter)
	r.Route("/backoffice-api", bapi.New(log, db).AttachChiRouter)

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
