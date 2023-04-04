package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"osoc-dialog/pkg/log"
)

func New(opts ...Option) *gin.Engine {
	o := options{
		logger:                 log.NewDiscardLogger(),
		docPath:                "undefined",
		middlewares:            []gin.HandlerFunc{},
		handleMethodNotAllowed: true,
		enableContextFallback:  true,
		pprof:                  false,
		pprofPrefix:            "debug/pprof",
	}
	for _, opt := range opts {
		opt(&o)
	}

	engine := gin.New()
	engine.HandleMethodNotAllowed = o.handleMethodNotAllowed
	engine.ContextWithFallback = o.enableContextFallback
	engine.Use(o.middlewares...)

	h := builtinHandlers{
		logger:          o.logger,
		docPath:         o.docPath,
		buildCommit:     o.buildCommit,
		buildTime:       o.buildTime,
		readinessProbes: o.readinessProbes,
	}

	engine.GET("/", h.root)
	engine.GET("/ready", h.readinessProbe)
	engine.GET("/live", h.livenessProbe)
	engine.GET("/doc", h.renderDoc)

	if o.pprof {
		pprof.Register(engine, o.pprofPrefix)
	}

	return engine
}
