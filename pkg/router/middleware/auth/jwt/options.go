package jwt

import (
	"osoc-dialog/pkg/log"
)

type Option func(*options)

type options struct {
	hmacSecret []byte
	logger     log.Logger
}

func HMACSecret(hmacSecret []byte) Option {
	return func(o *options) {
		o.hmacSecret = hmacSecret
	}
}

func Logger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}
