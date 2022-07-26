package application

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

type ViewTorrentInUILink string

type DownloadTorrentUsecase struct {
	fx.In
	Repo         TorrentRepository
	Downloader   TorrentDownloader
	Transmission TransmissionService
	Link         ViewTorrentInUILink
}

func (u *DownloadTorrentUsecase) Execute(torrentID int, notifier Notifier) error {
	torrent, err := u.Repo.FindByID(context.TODO(), torrentID)
	if err != nil {
		return err
	}

	filepath, err := u.Downloader.Download(torrent)
	if err != nil {
		return err
	}

	tTorrentID, err := u.Transmission.AddTorrent(filepath)
	if err != nil {
		return err
	}

	notifier.NotifyLink(
		fmt.Sprintf("Уже качаю '%s'", torrent.Name),
		&Link{Text: "Посмотреть", Url: string(u.Link)},
	)

	status := u.Transmission.WaitDone(tTorrentID)

	if status == TorrentDownloadTimeout {
		notifier.NotifyText(
			"Я больше не могу ждать, пока этот торрент скачается: %s", torrent.Name,
		)
		return fmt.Errorf("download torrent timeout: %d", torrent.ID)
	}

	if status == TorrentDownloaded {
		notifier.NotifyText("Скачалось! %s", torrent.Name)
		return nil
	}

	return fmt.Errorf("unknown download torrent status: %v", status)
}
