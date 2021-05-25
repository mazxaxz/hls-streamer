package videowatcher

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/radovskyb/watcher"

	"github.com/mazxaxz/hls-streamer/pkg/hls"
)

type watcherContext struct {
	hlsDir string
}

func New(hlsDirectory string) *watcherContext {
	c := watcherContext{
		hlsDir: hlsDirectory,
	}
	return &c
}

func (c *watcherContext) WatchCallback(ctx context.Context, op watcher.Op, fPath string, fi os.FileInfo) error {
	switch op {
	case watcher.Create, watcher.Move, watcher.Write:
		if fi.IsDir() {
			return nil
		}
		// we could use some kind of slugify package
		dir := strings.ReplaceAll(fi.Name(), ".", "_")
		options := hls.Options{
			VideoBitrate:      "5000k",
			VideoMaxRate:      "5350k",
			VideoBufSize:      "7500k",
			AudioBitrate:      "192k",
			AudioSamplingRate: 48000,
		}
		return hls.Transcode(ctx, fPath, path.Join(c.hlsDir, dir), options)
	default:
		return nil
	}
}
