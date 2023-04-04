package serviceprovider

import (
	"github.com/gin-gonic/gin"
	"osoc-dialog/internal/config"
	app "osoc-dialog/pkg/application"
	"osoc-dialog/pkg/log"
	"osoc-dialog/pkg/router"
	"osoc-dialog/pkg/router/middleware/logging"
	"osoc-dialog/pkg/router/middleware/recoverer"
	"osoc-dialog/pkg/router/middleware/servertiming"
)

func NewBaseRouter(conf *config.Config, logger log.Logger, version app.BuildVersion) *gin.Engine {
	return router.New(
		router.Logger(logger),
		router.DocPath(conf.App.SwaggerFolder),
		router.BuildCommit(version.Commit),
		router.BuildTime(version.Time),
		router.Middlewares(
			recoverer.New(
				recoverer.Logger(logger),
			),
			servertiming.New(),
			logging.New(
				logging.Level(conf.App.LogLevel),
				logging.Env(conf.App.Environment),
				logging.Fallback(logger),
				logging.Prettify(conf.App.PrettyLogs),
			),
		))
}
