package config

// PromConfig - конфиг для прометея
type PromConfig struct {
	Port    string `envconfig:"PROM_PORT" default:"9100"`
	Handle  string `envconfig:"PROM_HANDLE" default:"/metrics"`
	GuiPort string `envconfig:"PROM_GUI_PORT" default:"9090"`
}

func GetPrometheusConfig(config *Config) PromConfig {
	return config.Prometheus
}
