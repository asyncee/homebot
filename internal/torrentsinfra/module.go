package torrentsinfra

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewRutrackerTorrentRepository,
		NewDownloader,
		NewTransmissionService,
	),
)
