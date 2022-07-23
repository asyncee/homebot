package torrentsinfra

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/asyncee/homebot/internal/torrents/application"
	"github.com/asyncee/homebot/internal/torrents/domain"
	"github.com/asyncee/homebot/pkg/rutracker"
)

type DownloadTorrentsDir string

type downloader struct {
	client *rutracker.RutrackerClient
	dir    DownloadTorrentsDir
}

func NewDownloader(client *rutracker.RutrackerClient, dir DownloadTorrentsDir) application.TorrentDownloader {
	return &downloader{client: client, dir: dir}
}

func (d *downloader) Download(torrent *domain.Torrent) (application.Filepath, error) {
	fetchedFile, err := d.client.FetchTorrent(torrent.DownloadUrl, torrent.DownloadCsrfToken)
	if err != nil {
		return "", fmt.Errorf("failed to fetch torrent file %d from %s: %s", torrent.ID, torrent.DownloadUrl, err)
	}

	err = os.MkdirAll(string(d.dir), 0777)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s for torrent %d: %s", d.dir, torrent.ID, err)
	}
	fpath := filepath.Join(string(d.dir), fetchedFile.Name)
	f, err := os.Create(fpath)
	if err != nil {
		return "", fmt.Errorf("failed to create torrent file %s on filesystem (torrent id %d)", fpath, torrent.ID)
	}

	_, err = io.Copy(f, fetchedFile.Body)
	if err != nil {
		return "", fmt.Errorf("failed to write torrent file %s (torrent id %d)", fpath, torrent.ID)
	}

	return application.Filepath(fpath), nil
}
