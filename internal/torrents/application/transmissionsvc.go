package application

type TransmissionTorrentId int64

type DownloadTorrentStatus int

const (
	TorrentDownloadTimeout DownloadTorrentStatus = 0
	TorrentDownloaded      DownloadTorrentStatus = 1
	TorrentStatusUnknown   DownloadTorrentStatus = 2
)

type TransmissionService interface {
	AddTorrent(filepath Filepath) (TransmissionTorrentId, error)
	WaitDone(torrentID TransmissionTorrentId) DownloadTorrentStatus
}
