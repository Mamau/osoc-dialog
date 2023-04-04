package router

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"osoc-dialog/pkg/healthcheck"
	"osoc-dialog/pkg/router/response"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"osoc-dialog/pkg/log"
)

type builtinHandlers struct {
	logger          log.Logger
	docPath         string
	buildCommit     string
	buildTime       time.Time
	readinessProbes []healthcheck.ProbeFunc
}

func (h *builtinHandlers) readinessProbe(c *gin.Context) {
	if len(h.readinessProbes) == 0 {
		response.Render(c, http.StatusOK, "Service Ready")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	statuses, err := healthcheck.CheckProbes(ctx, h.readinessProbes)
	if err != nil {
		h.logger.Err(err).Msg("Could not perform readiness check")
		response.Render(c, http.StatusServiceUnavailable, "Service Not Ready")
		return
	}

	result := make(map[string]*healthcheck.ProbeStatus, len(statuses))
	ready := true

	for _, v := range statuses {
		result[v.Component] = &v
		if v.Critical && !v.Ready {
			ready = false
		}
	}

	if ready {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusServiceUnavailable)
	}
	_ = render.IndentedJSON{Data: map[string]any{"data": result}}.Render(c.Writer)
}

func (h *builtinHandlers) livenessProbe(c *gin.Context) {
	response.Render(c, http.StatusOK, "Service Alive")
}

func (h *builtinHandlers) printVersion(c *gin.Context) {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	response.Render(c, http.StatusOK, struct {
		Hostname string `json:"hostname"`
		Commit   string `json:"commit"`
		Time     string `json:"time"`
	}{
		Hostname: hostname,
		Commit:   h.buildCommit,
		Time:     h.buildTime.In(time.UTC).Format(time.RFC3339),
	})
}

func (h *builtinHandlers) renderDoc(c *gin.Context) {
	filename := filepath.Join(h.docPath, "api.swagger.json")

	f, err := os.Open(filename)
	if err != nil {
		h.logger.Err(err).Msgf("could not open file: %v", filename)
		response.Render(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	b, err := io.ReadAll(f)
	if err != nil {
		h.logger.Err(err).Msgf("could not read file: %v", filename)
		response.Render(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	mime := mimetype.Detect(b)
	c.Writer.Header().Add("Content-Type", mime.String())

	if _, err := c.Writer.Write(b); err != nil {
		h.logger.Err(err).Msgf("could not write data")
		response.Render(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func (h *builtinHandlers) root(c *gin.Context) {
	response.Render(c, http.StatusOK, http.StatusText(http.StatusOK))
}

func (h *builtinHandlers) notFound(c *gin.Context) {
	response.Render(c, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}
