package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/mazxaxz/hls-streamer/cmd/apid/config"
	"github.com/mazxaxz/hls-streamer/cmd/apid/hlshttphandler"
	"github.com/mazxaxz/hls-streamer/internal/hls"
	"github.com/mazxaxz/hls-streamer/pkg/logger"
	"github.com/mazxaxz/hls-streamer/pkg/osdir"
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
	hlsService := hls.NewService(cfg.HLSDirectory)
	hlsHttpHandler := hlshttphandler.New(cfg.HLSDirectory, log)

	transcodeExisting(ctx, cfg.RawDirectory, cfg.HLSDirectory)
	go osdir.Watch(ctx, log, cfg.RawDirectory, hlsService.WatchCallback)

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

func transcodeExisting(ctx context.Context, rawDir, hlsDir string) {
	if _, err := os.Stat(rawDir); os.IsNotExist(err) {
		if err := os.MkdirAll(rawDir, os.ModePerm); err != nil {
			log.WithContext(ctx).Error(err)
			return
		}
	}
	files, err := ioutil.ReadDir(rawDir)
	if err != nil {
		log.WithContext(ctx).Error(err)
		return
	}
	for _, file := range files {
		select {
		case <-ctx.Done():
			return
		default:
			if file.IsDir() {
				continue
			}
			log.WithContext(ctx).WithField("file", file.Name()).Info("transcoding...")
			dir := strings.ReplaceAll(file.Name(), ".", "_")
			if err := hls.Transcode(ctx, path.Join(rawDir, file.Name()), path.Join(hlsDir, dir)); err != nil {
				log.WithContext(ctx).Error(err)
			}
			log.WithContext(ctx).WithField("file", file.Name()).Info("transcoding finished")
		}
	}
}
