//go:build wireinject

package main

import (
	"osoc-dialog/internal/serviceprovider"

	"github.com/google/wire"

	"osoc-dialog/pkg/application"
)

func newApp() (*application.App, func(), error) {
	panic(wire.Build(
		serviceprovider.ProviderSet,
		createApp,
	))
}
