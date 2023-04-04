package recoverer

import (
	"osoc-dialog/pkg/log"
)

type Option func(o *options)

type options struct {
	logger log.Logger
}

func Logger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}
