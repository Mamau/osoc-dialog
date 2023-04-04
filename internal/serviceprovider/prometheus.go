package serviceprovider

import (
	"osoc-dialog/internal/config"

	"osoc-dialog/pkg/log"
	"osoc-dialog/pkg/transport/prom"
)

func NewPrometheus(config config.PromConfig, logger log.Logger) *prom.Server {
	server := prom.New(
		prom.Logger(logger),
		prom.GuiPort(config.GuiPort),
		prom.Port(config.Port),
		prom.Handle(config.Handle),
	)

	return server
}
