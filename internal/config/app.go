package config

import "time"

// App конфиг приложения
type App struct {
	Name          string `envconfig:"APP_NAME" default:"osoc-dialog"`
	Environment   string `envconfig:"APP_ENV" default:"prod"`
	Host          string `envconfig:"APP_HOST" default:"localhost"`
	Port          string `envconfig:"APP_PORT" default:"8081"`
	LogLevel      string `envconfig:"APP_LOG_LEVEL" default:"info"`
	PrettyLogs    bool   `envconfig:"APP_LOG_PRETTY" default:"false"`
	SwaggerFolder string `envconfig:"SWAGGER_FOLDER" default:"swagger"`
	TZ            string `envconfig:"TZ" default:"Europe/Moscow"`
	GinMode       string `envconfig:"GIN_MODE" default:"release"`

	RateLimitInterval time.Duration `envconfig:"RATE_LIMIT_INTERVAL" default:"10ms"`
	RateLimitBurst    int           `envconfig:"RATE_LIMIT_BURST" default:"200"`
	DaemonRunInterval time.Duration `envconfig:"DAEMON_RUN_INTERVAL" default:"5s"`
}

func GetAppConfig(config *Config) App {
	return config.App
}
