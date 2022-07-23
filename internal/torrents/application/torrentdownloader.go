package application

import "github.com/asyncee/homebot/internal/torrents/domain"

type Filepath string

type TorrentDownloader interface {
	Download(torrent *domain.Torrent) (Filepath, error)
}
