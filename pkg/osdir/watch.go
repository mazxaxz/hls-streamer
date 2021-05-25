package osdir

import (
	"context"
	"os"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/sirupsen/logrus"
)

type WatchFnCallback func(ctx context.Context, op watcher.Op, path string, fi os.FileInfo) error

func Watch(ctx context.Context, log *logrus.Logger, path string, cb WatchFnCallback) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.WithContext(ctx).Error(err)
			return
		}
	}

	w := watcher.New()
	defer w.Close()

	if err := w.Add(path); err != nil {
		log.WithContext(ctx).WithField("path", path).Error(err)
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-w.Closed:
				return
			case err := <-w.Error:
				log.Error(err)
			case event := <-w.Event:
				if err := cb(ctx, event.Op, event.Path, event.FileInfo); err != nil {
					log.WithField("event", event.String()).WithContext(ctx).Error(err)
				}
			}
		}
	}()

	if err := w.Start(1 * time.Second); err != nil {
		log.WithContext(ctx).Error(err)
	}
}
