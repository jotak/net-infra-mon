package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jotak/net-infra-mon/pkg/config"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "server")

func Start(ctx context.Context, cfg *config.Config) {
	router := setupRoutes(ctx, cfg)

	httpServer := Default(&http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
	})

	if cfg.Server.CertPath != "" && cfg.Server.KeyPath != "" {
		log.Infof("listening on https://:%d", cfg.Server.Port)
		panic(httpServer.ListenAndServeTLS(cfg.Server.CertPath, cfg.Server.KeyPath))
	}
	log.Infof("listening on http://:%d", cfg.Server.Port)
	panic(httpServer.ListenAndServe())
}
