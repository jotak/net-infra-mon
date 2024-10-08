package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/jotak/net-infra-mon/pkg/config"
	"github.com/jotak/net-infra-mon/pkg/jobs"
	"github.com/jotak/net-infra-mon/pkg/server"
)

var (
	buildVersion = "unknown"
	buildDate    = "unknown"
	app          = "net-infra-mon"
	configPath   = flag.String("config", "", "path to the configuration file")
	versionFlag  = flag.Bool("v", false, "print version")
	log          = logrus.WithField("module", "main")
)

func main() {
	flag.Parse()

	appVersion := fmt.Sprintf("%s [build version: %s, build date: %s]", app, buildVersion, buildDate)
	if *versionFlag {
		fmt.Println(appVersion)
		os.Exit(0)
	}
	log.Infof("Starting %s", appVersion)

	cfg, err := config.Read(*configPath)
	if err != nil {
		log.WithError(err).Fatal("error reading config file")
	}
	log.Infof("Configuration: %+v", cfg)

	lvl, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Errorf("Log level %s not recognized", cfg.LogLevel)
	} else {
		logrus.SetLevel(lvl)
		log.Infof("Log level set to %s", cfg.LogLevel)
	}

	// Pprof
	if cfg.PProfPort > 0 {
		go func() {
			log.WithField("port", cfg.PProfPort).Info("starting PProf HTTP listener")
			log.WithError(http.ListenAndServe(fmt.Sprintf(":%d", cfg.PProfPort), nil)).
				Error("PProf HTTP listener stopped working")
		}()
	}

	go func() {
		jobs.Run(context.Background())
	}()

	// Blocking call
	server.Start(context.Background(), cfg)
}
