package application

import (
	"context"

	"github.com/asyncee/homebot/internal/torrents/domain"
)

type TorrentRepository interface {
	FindByName(ctx context.Context, name string) ([]domain.Torrent, error)
}
