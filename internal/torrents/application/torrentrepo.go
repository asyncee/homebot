package application

import (
	"context"

	"github.com/asyncee/homebot/internal/torrents/domain"
)

type TorrentRepository interface {
	FindByName(ctx context.Context, name string) ([]domain.Torrent, error)
	FindByID(ctx context.Context, id int) (*domain.Torrent, error)
}
