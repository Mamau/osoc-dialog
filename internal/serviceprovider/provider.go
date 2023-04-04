package serviceprovider

import (
	"github.com/google/wire"
	v1 "osoc-dialog/internal/api/http/v1"
	"osoc-dialog/internal/config"
	dr "osoc-dialog/internal/repository/dialog"
	"osoc-dialog/internal/usecase/dialog"
	"osoc-dialog/pkg/application"
)

var ProviderSet = wire.NewSet(
	wire.Bind(new(v1.DialogProvider), new(*dialog.Service)), dialog.NewService,
	wire.Bind(new(dialog.MessageStorage), new(*dr.Repository)), dr.New,
	NewBaseRouter,
	application.GetBuildVersion,
	config.GetPrometheusConfig,
	config.GetConfig,
	config.GetAppConfig,
	config.GetMysqlConfig,
	config.GetProxyMysqlConfig,
	NewHttp,
	NewMysql,
	NewProxyMysql,
	NewPrometheus,
	NewLogger,
	v1.NewRouter,
)
