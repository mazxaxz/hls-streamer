package hlshttphandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/mazxaxz/hls-streamer/pkg/rest"
)

type handlerContext struct {
	logger *logrus.Logger
}

func New(l *logrus.Logger) rest.SetupRouterer {
	c := handlerContext{
		logger: l,
	}
	return &c
}

func (c *handlerContext) SetupRouter(r *gin.RouterGroup) {
	r.GET("/hls", c.serveHLS)
}

func (c *handlerContext) serveHLS(cGin *gin.Context) {
	cGin.String(200, "test")
}
