package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"osoc-dialog/pkg/log"
)

func NewRouter(
	engine *gin.Engine,
	logger log.Logger,
	dp DialogProvider,
) http.Handler {
	commonGroup := engine.Group("/api/v1")
	commonGroup.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })

	dialogGroup := commonGroup.Group("/dialog")
	{
		newDialogRoutes(dialogGroup, logger, dp)
	}

	return engine
}
