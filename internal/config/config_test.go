package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	if err := os.Setenv("APP_NAME", "test-osoc-dialog"); err != nil {
		assert.Nil(t, err)
	}

	conf, _ := GetConfig()
	assert.Equal(t, "test-osoc-dialog", conf.App.Name)

	os.Clearenv()
	conf, _ = GetConfig()
	assert.Equal(t, "osoc-dialog", conf.App.Name)
}

func TestGetMysqlConfig(t *testing.T) {
	conf, _ := GetConfig()
	mc := GetMysqlConfig(conf)
	assert.IsType(t, Mysql{}, mc)
}

func TestGetPrometheusConfig(t *testing.T) {
	conf, _ := GetConfig()
	pc := GetPrometheusConfig(conf)
	assert.IsType(t, PromConfig{}, pc)
}

func TestGetRedisConfig(t *testing.T) {
	conf, _ := GetConfig()
	pc := GetRedisConfig(conf)
	assert.IsType(t, Redis{}, pc)
}
