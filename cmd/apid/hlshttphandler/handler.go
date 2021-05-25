package hlshttphandler

import (
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mazxaxz/hls-streamer/pkg/rest"
)

type handlerContext struct {
	hlsDir string
	logger *logrus.Logger
}

func New(hlsDirectory string, l *logrus.Logger) rest.SetupRouterer {
	c := handlerContext{
		hlsDir: hlsDirectory,
		logger: l,
	}
	return &c
}

func (c *handlerContext) SetupRouter(r *gin.RouterGroup) {
	r.GET("/:video/:file", c.serveHLS)
}

func (c *handlerContext) serveHLS(cGin *gin.Context) {
	video := cGin.Param("video")
	if video == "" {
		cGin.AbortWithStatusJSON(http.StatusBadRequest, "video parameter cannot be empty")
		return
	}
	file := cGin.Param("file")
	if file == "" {
		cGin.AbortWithStatusJSON(http.StatusBadRequest, "file parameter cannot be empty")
		return
	}

	videoDir := path.Join(c.hlsDir, video)
	cGin.FileFromFS(file, http.Dir(videoDir))
}
