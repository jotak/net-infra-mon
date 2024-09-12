package server

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/jotak/net-infra-mon/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func setupRoutes(ctx context.Context, cfg *config.Config) *mux.Router {
	r := mux.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	return r
}
