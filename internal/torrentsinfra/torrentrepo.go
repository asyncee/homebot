package torrentsinfra

import (
	"context"
	"fmt"

	"github.com/asyncee/homebot/internal/torrents/application"
	"github.com/asyncee/homebot/internal/torrents/domain"
	"github.com/asyncee/homebot/pkg/logging"
	"github.com/asyncee/homebot/pkg/rutracker"
	"github.com/golang/groupcache/lru"
)

func NewRutrackerTorrentRepository(
	client *rutracker.RutrackerClient,
	logger logging.Logger,
) application.TorrentRepository {
	return &RutrackerTorrentsRepository{
		client: client,
		logger: logger,
		cache:  lru.New(1000),
	}
}

type RutrackerTorrentsRepository struct {
	client *rutracker.RutrackerClient
	cache  *lru.Cache
	logger logging.Logger
}

func (re *RutrackerTorrentsRepository) FindByName(ctx context.Context, name string) ([]domain.Torrent, error) {
	torrents, err := re.client.Search(name)
	if err != nil {
		return nil, err
	}

	if len(torrents) == 0 {
		return []domain.Torrent{}, nil
	}

	results := make([]domain.Torrent, 0, len(torrents))

	for i := range torrents {
		torrent := re.toDomainModel(torrents[i])
		re.cache.Add(torrent.ID, torrent)
		results = append(results, torrent)
	}

	return results, nil
}

func (re *RutrackerTorrentsRepository) FindByID(ctx context.Context, id int) (*domain.Torrent, error) {
	if torrent, ok := re.cache.Get(id); ok {
		return torrent.(*domain.Torrent), nil
	}
	return nil, fmt.Errorf("can't find torrent by id: key %d not found in cache", id)
}

func (re *RutrackerTorrentsRepository) toDomainModel(item rutracker.Torrent) domain.Torrent {
	return domain.Torrent{
		ID:                item.ID,
		Status:            item.Status,
		Name:              item.Name,
		Size:              item.Size,
		Seeders:           item.Seeders,
		DownloadUrl:       item.DownloadUrl,
		DownloadCsrfToken: item.DownloadToken,
		URL:               item.URL,
		Category:          item.Category,
	}

}
