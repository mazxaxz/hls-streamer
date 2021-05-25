package hls

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/radovskyb/watcher"
)

func (c *Service) WatchCallback(ctx context.Context, op watcher.Op, fPath string, fi os.FileInfo) error {
	switch op {
	case watcher.Create, watcher.Move, watcher.Write:
		if fi.IsDir() {
			return nil
		}

		// we could use some kind of slugify package
		dir := strings.ReplaceAll(fi.Name(), ".", "_")
		return Transcode(ctx, fPath, path.Join(c.hlsDir, dir))
	default:
		return nil
	}
}
