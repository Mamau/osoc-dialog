package mysql

import (
	"time"
)

type Option func(*options)

type options struct {
	host                 string
	port                 int
	user                 string
	password             string
	dbName               string
	parseTime            bool
	maxOpenConnections   int
	maxIdleConnections   int
	connMaxLifetime      time.Duration
	verificationRequired bool
	verificationTimeout  time.Duration
}

func Host(host string) Option {
	return func(c *options) {
		c.host = host
	}
}

func Port(port int) Option {
	return func(c *options) {
		c.port = port
	}
}

func User(user string) Option {
	return func(c *options) {
		c.user = user
	}
}

func Password(password string) Option {
	return func(c *options) {
		c.password = password
	}
}

func DBName(dbName string) Option {
	return func(c *options) {
		c.dbName = dbName
	}
}

func ParseTime(parseTime bool) Option {
	return func(c *options) {
		c.parseTime = parseTime
	}
}

func MaxOpenConnections(connections int) Option {
	return func(c *options) {
		c.maxOpenConnections = connections
	}
}

func MaxIdleConnections(connections int) Option {
	return func(c *options) {
		c.maxIdleConnections = connections
	}
}

func ConnMaxLifetime(cml time.Duration) Option {
	return func(c *options) {
		c.connMaxLifetime = cml
	}
}

func VerificationRequired(require bool) Option {
	return func(c *options) {
		c.verificationRequired = require
	}
}

func VerificationTimeout(timeout time.Duration) Option {
	return func(c *options) {
		c.verificationTimeout = timeout
	}
}
