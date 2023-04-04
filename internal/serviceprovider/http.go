package serviceprovider

import (
	nh "net/http"

	"osoc-dialog/internal/config"

	"osoc-dialog/pkg/log"
	"osoc-dialog/pkg/transport/http"
)

func NewHttp(handler nh.Handler, conf *config.Config, logger log.Logger) *http.Server {
	return http.New(
		http.Logger(logger),
		http.Handler(handler),
		http.Addr(conf.App.Port),
	)
}
