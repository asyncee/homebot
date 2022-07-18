package transmission

import (
	"context"
	"fmt"

	"github.com/hekmon/transmissionrpc/v2"
)

type TransmissionService struct {
	client *transmissionrpc.Client
}

func New(
	host string,
	user string,
	password string,
) (*TransmissionService, error) {
	client, err := transmissionrpc.New(host, user, password, nil)
	if err != nil {
		return nil, err
	}
	return &TransmissionService{
		client: client,
	}, nil
}

func (s *TransmissionService) TestConnection() error {
	ctx, cancel := context.WithCancel(context.TODO())
	ok, serverVersion, serverMinimumVersion, err := s.client.RPCVersion(ctx)
	defer cancel()
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf(
			"remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion,
			transmissionrpc.RPCVersion,
			serverMinimumVersion,
		)
	}
	return nil
}

func (s *TransmissionService) ListTorrents() ([]transmissionrpc.Torrent, error) {
	torrents, err := s.client.TorrentGetAll(context.TODO())
	if err != nil {
		return nil, err
	}
	return torrents, nil
}

type DownloadedTorrent struct {
	ID   int64
	Name string
}

func (s *TransmissionService) DownloadTorrent(filepath string) (*DownloadedTorrent, error) {
	torrent, err := s.client.TorrentAddFile(context.TODO(), filepath)
	if err != nil {
		return nil, err
	}
	return &DownloadedTorrent{
		ID:   *torrent.ID,
		Name: *torrent.Name,
	}, nil
}

type TorrentStatus int

const (
	TorrentDone    = 1
	TorrentRunning = 2
	TorrentMissing = 3
	TorrentError   = 4
)

func (s *TransmissionService) TorrentIsDone(torrentId int64) (TorrentStatus, error) {
	torrents, err := s.client.TorrentGet(context.TODO(), []string{"status"}, []int64{torrentId})
	if err != nil {
		return TorrentError, err
	}

	if len(torrents) == 0 {
		return TorrentMissing, nil
	}

	status := *torrents[0].Status
	if status == transmissionrpc.TorrentStatusSeed || status == transmissionrpc.TorrentStatusSeedWait || status == transmissionrpc.TorrentStatusStopped {
		return TorrentDone, nil
	}

	return TorrentRunning, nil
}
