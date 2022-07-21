package torrentsinfra

import (
	"context"
	"strconv"

	"github.com/asyncee/homebot/internal/torrents/application"
	"github.com/asyncee/homebot/internal/torrents/domain"
	"github.com/asyncee/homebot/pkg/rutracker"
)

func NewRutrackerTorrentRepository(client *rutracker.RutrackerClient) application.TorrentRepository {
	return &RutrackerTorrentsRepository{client: client}
}

type RutrackerTorrentsRepository struct {
	client *rutracker.RutrackerClient
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
		results = append(results, re.toDomainModel(torrents[i]))
	}

	return results, nil
}

func (re *RutrackerTorrentsRepository) toDomainModel(item rutracker.Torrent) domain.Torrent {
	return domain.Torrent{
		ID:                strconv.Itoa(item.ID),
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
