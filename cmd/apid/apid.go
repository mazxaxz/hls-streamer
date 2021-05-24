package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mazxaxz/hls-streamer/cmd/apid/config"
	"github.com/mazxaxz/hls-streamer/cmd/apid/hlshttphandler"
	"github.com/mazxaxz/hls-streamer/pkg/logger"
	"github.com/mazxaxz/hls-streamer/pkg/shutdown"
)

var log = logrus.New()

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.Configure(log, cfg.Logger); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// handlers
	hlsHttpHandler := hlshttphandler.New(log)

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      setupRouting(hlsHttpHandler),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func(ctx context.Context, srv http.Server) {
		log.Info(fmt.Sprintf("Starting server on port: %s", srv.Addr))
		if err := srv.ListenAndServe(); err != nil {
			log.Info("Closing server...")
		}
	}(ctx, srv)

	shutdown.Wait(cancel, log, srv)
}
