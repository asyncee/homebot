package torrentsinfra

import "go.uber.org/fx"

var Module = fx.Provide(NewRutrackerTorrentRepository)
