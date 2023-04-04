package config

type Mysql struct {
	Host      string `envconfig:"MY_HOST" default:"mysql"`
	Port      int    `envconfig:"MY_PORT" default:"3306"`
	User      string `envconfig:"MY_USER" default:"root"`
	Password  string `envconfig:"MY_PASSWORD" default:"root"`
	DbName    string `envconfig:"MY_DB_NAME" default:"osoc-dialog"`
	ParseTime bool   `envconfig:"MY_PARSE_TIME" default:"true"`
}

type ProxyMysql struct {
	Host      string `envconfig:"MY_PROXY_HOST" default:"proxysql"`
	Port      int    `envconfig:"MY_PROXY_PORT" default:"6033"`
	User      string `envconfig:"MY_PROXY_USER" default:"test"`
	Password  string `envconfig:"MY_PROXY_PASSWORD" default:"pzjqUkMnc7vfNHET"`
	DbName    string `envconfig:"MY_PROXY_DB_NAME" default:"test"`
	ParseTime bool   `envconfig:"MY_PROXY_PARSE_TIME" default:"true"`
}

func GetMysqlConfig(config *Config) Mysql {
	return config.Mysql
}

func GetProxyMysqlConfig(config *Config) ProxyMysql {
	return config.ProxyMysql
}
