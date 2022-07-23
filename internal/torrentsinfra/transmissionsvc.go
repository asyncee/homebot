package torrentsinfra

import (
	"time"

	"github.com/asyncee/homebot/internal/torrents/application"
	"github.com/asyncee/homebot/pkg/transmission"
)

type transmissionsvc struct {
	client             *transmission.TransmissionClient
	pollStatusDuration time.Duration
	pollStatusTimeout  time.Duration
}

type PollStatusDuration time.Duration
type PollStatusTimeout time.Duration

func NewTransmissionService(
	client *transmission.TransmissionClient,
	pollStatusDuration PollStatusDuration,
	pollStatusTimeout PollStatusTimeout,
) application.TransmissionService {
	return &transmissionsvc{
		client:             client,
		pollStatusDuration: time.Duration(pollStatusDuration),
		pollStatusTimeout:  time.Duration(pollStatusTimeout),
	}
}

func (tr *transmissionsvc) AddTorrent(filepath application.Filepath) (application.TransmissionTorrentId, error) {
	torrent, err := tr.client.DownloadTorrent(string(filepath))
	if err != nil {
		return -1, err
	}
	return application.TransmissionTorrentId(torrent.ID), nil
}

func (tr *transmissionsvc) WaitDone(torrentID application.TransmissionTorrentId) application.DownloadTorrentStatus {
	ticker := time.NewTicker(tr.pollStatusDuration)
	timer := time.NewTimer(tr.pollStatusTimeout)

	for range ticker.C {
		select {
		case <-timer.C:
			ticker.Stop()
			return application.TorrentDownloadTimeout
		case <-ticker.C:
			// Periodically check if torrent is downloaded
			status, err := tr.client.TorrentIsDone(int64(torrentID))
			if err != nil {
				continue
			}

			if status == transmission.TorrentDone {
				if !timer.Stop() {
					<-timer.C
				}
				ticker.Stop()
				return application.TorrentDownloaded
			}

			if status == transmission.TorrentRunning {
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(tr.pollStatusTimeout)
			}
		}
	}

	return application.TorrentStatusUnknown
}
